// destructure the L402 headers from the URL
function destructL402Authenticate(header) {
    // strip off "L402 " from the beginning of the header
    const l402Header = header.slice(5);
    let [token, invoice] = l402Header.split(", ")

    // strip off "token=" and "invoice=" from the beginning of the token and invoice
    token = token.slice(6);
    invoice = invoice.slice(8);

    return [token, invoice];
  }

  async function fetchWithL402(url, fetchArgs, headerElements) {
    const { l402ElementId, authorizationElementId } = headerElements;
    let res = await fetch(url, fetchArgs);
    if (res.status === 402) {
      const [token, invoice] = await destructL402Authenticate(res.headers.get("WWW-Authenticate"));

      // Display L402 headers regardless of whether webln is defined or not.
      displayL402Header(l402ElementId, token, invoice, false);

      if (typeof window.webln !== "undefined") {
        await window.webln.enable();
        const { preimage } = await window.webln.sendPayment(invoice);
        if (!!preimage) {
          let authorizationValue = `L402 ${token}:${preimage}`;
          fetchArgs.headers["Authorization"] = authorizationValue;

          // Display Authorization header
          displayL402Header(authorizationElementId, token, preimage, true);

          res = await fetch(url, fetchArgs);
          if (res.status === 402) {
            alert("Payment failed");
          } else if (res.status === 200) {
            return res;
          }
        } else {
          alert("Payment failed");
        }
      } else {
        throw new Error("Payment required. L402 header displayed.");
      }
    } else if (res.status === 200) {
      const data = await res.json();
      return data;
    }
    return res;
  }



  function displayL402Header(elementId, token, invoice, auth) {
    if (auth) {
      document.querySelector("#" + elementId).textContent = `Authorization: L402 ${token}:${invoice}`;
      return;
    } else if (invoice) {
      document.querySelector("#" + elementId).textContent = `WWW-Authenticate: L402 token=${token}, invoice=${invoice}`;
      return;
    }
  }


  async function generateChatCompletion() {
    if (typeof window.webln === "undefined") {
        alert("This demo of Matador requires a bitcoin browser extension, you can download one at https://getalby.com");
        return;
    }
    const model = document.querySelector("#model").value;
    const systemRole = document.querySelector("#systemRole").value;
    const userRole = document.querySelector("#userRole").value;
    const payload = {
      model: model,
      messages: [
        { role: "system", content: systemRole },
        { role: "user", content: userRole },
      ],
    };

    // Replace the URL with your Matador server's URL
    const response = await fetchWithL402("https://matador-ai.replit.app/v1/chat/completions", {
      method: "POST",
      body: JSON.stringify(payload),
      mode: "cors",
      headers: { "Content-Type": "application/json", Accept: "application/json" },
    }, { l402ElementId: "chatCompletionL402Header", authorizationElementId: "chatCompletionAuthorizationHeader" });


    // Displaying the response
    const json = await response.json();
    document.querySelector("#chatCompletionResponse").textContent = JSON.stringify(json, null, 2);
  }

  async function generateImage() {
      if (typeof window.webln === "undefined") {
        alert("This demo of Matador requires a bitcoin browser extension, you can download one at https://getalby.com");
        return;
    }
    const prompt = document.querySelector("#prompt").value;
    const payload = {
      prompt: prompt,
      n: 1,
      size: "512x512",
      response_format: "url",
    };

    // Replace the URL with your Matador server's URL
    const response = await fetchWithL402("https://matador-ai.replit.app/v1/images/generations", {
      method: "POST",
      body: JSON.stringify(payload),
      mode: "cors",
      headers: { "Content-Type": "application/json", Accept: "application/json" },
    }, { l402ElementId: "imageGenerationL402Header", authorizationElementId: "imageGenerationAuthorizationHeader" });


    // Displaying the response
    const json = await response.json();
    document.querySelector("#imageGenerationResponse").textContent = JSON.stringify(json, null, 2);
  }

  async function generateEmbedding() {
      if (typeof window.webln === "undefined") {
        alert("This demo of Matador requires a bitcoin browser extension, you can download one at https://getalby.com");
        return;
    }
    const input = document.querySelector("#embeddingInput").value;
    const model = document.querySelector("#embeddingModel").value;
    const payload = {
      input: input,
      model: model
    };

    const response = await fetchWithL402("http://localhost:8080/v1/embeddings", {
      method: "POST",
      body: JSON.stringify(payload),
      mode: "cors",
      headers: { "Content-Type": "application/json", Accept: "application/json" },
    }, { l402ElementId: "embeddingL402Header", authorizationElementId: "embeddingAuthorizationHeader" });

    // Displaying the response
    const json = await response.json();
    // limiting the number of embeddings displayed
    json.data[0].embedding = json.data[0].embedding.slice(0, 10).concat(["... (total 1536 values)"]);
    document.querySelector("#embeddingResponse").textContent = JSON.stringify(json, null, 2);
  }

  function expandText() {
    let content = document.getElementById('expandable');
    content.style.webkitLineClamp = 'unset';
  }

  function toggleSwitch(element) {
    const isChecked = element.checked;
    const chatCompletionSection = document.querySelector('.chat-completion');
    const imageGenerationSection = document.querySelector('.image-generation');
    if (isChecked) {
      chatCompletionSection.style.display = 'none';
      imageGenerationSection.style.display = 'block';
    } else {
      chatCompletionSection.style.display = 'block';
      imageGenerationSection.style.display = 'none';
    }
  }