// destructure the L402 headers from the URL
function destructL402Authenticate(header) {
    // strip off "L402 " from the beginning of the header
    const l402Header = header.slice(5);
    let [token, invoice] = l402Header.split(", ")
    console.log("token:", token)
    console.log("invoice:", invoice)
    // strip off "token=" and "invoice=" from the beginning of the token and invoice
    token = token.slice(6).replace("\"", "");
    invoice = invoice.slice(8).replace("\"", "");
    console.log("token:", token)
    console.log("invoice:", invoice)

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
