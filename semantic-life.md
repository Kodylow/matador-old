# API Usage Guide

FOR TESTING BY DEVELOPERS ONLY. CANNOT BE USED FOR COMMERCIAL PURPOSES. PROVIDED UNDER MIT LICENSE. USE AT YOUR OWN RISK. DO NOT REPLICATE, DO NOT SELL. FOR TESTING, RESEARCH, AND DEMONSTRATION ONLY. NO PRIVATE DATA. SOME DATA RENDERED FROM OPENAI API, SOME BY HUMANS, SOME BY OTHER ENGINES. ALL DERIVATIVE TERMS APPLY. 

CONTACT TWITTER.COM/SEANMCDONALDXYZ TO INQUIRE ABOUT COMMERCIAL API.

Can be used to get a rich prompt for api calls. 

Search function: find agents by keyword or job title. Use agent name as id to get full record.

Returns: 
"name": "John Doe",
  "job": "Agent",
  "name": Agent name.
  "job": Agent job title.
  "field": The field the agent works in.
  "character_trait": Character traits.
  "character_history": Backstory.
  "demo": Demographic information.
  "fren": Friend network.
  "persona": Likely persona details.
  "goal": Goals based on other details
  "role": (based on config file generating agent see our Brancher repository on twitter.com/semanticlife)
  "id": (based on config file)

Request an API Key at semantic-life.com. 

# 1. Find the Agents You Want by Name/Job

### A. Retrieve all agents
Endpoint: GET /api/v1/resources/agents

curl -X GET -H "X-User-ID: <user_id>" -H "X-User-Key: <user_key>" /api/v1/resources/agents


### B. Get Name / Job by Keyword Search of ONLY Jobs

curl -X GET -H "X-User-ID: {}" -H "X-User-Key: {}" /api/v1/resources/agents?job=investment

Response:
Status code: 200 OK
Response body:

[
  {
    "name": "Agent1",
    "job": "Investment Banker"
  },
  {
    "name": "Agent2",
    "job": "Investment Analyst"
  }
]

### C. Get Name / Job by Keyword Search of All Fields

curl -X GET -H "X-User-ID: {}" -H "X-User-Key: {}" /api/v1/resources/agents?keyword=sports%20fan
Response:
Status code: 200 OK
Response body:


Copy

Insert
[
  {
    "name": "Agent1",
    "job": "Sports Fan"
  },
  ...
]


# 2. Get Full Record using Name as id

curl -X GET -H "X-User-ID: {}" -H "X-User-Key: {}" /api/v1/resources/agents?name=John%20Doe

Response:
Status code: 200 OK
Response body:

[
{
  "name": "John Doe",
  "job": "Agent",
  "name": Agent name.
  "job": Agent job title.
  "field": The field the agent works in.
  "character_trait": Character traits.
  "character_history": Backstory.
  "demo": Demographic information.
  "fren": Friend network.
  "persona": Likely persona details.
  "goal": Goals based on other details
  "role": (based on config file generating agent see our Brancher repository on twitter.com/semanticlife)
  "id": (based on config file)
}
]



Please note that you need to replace <user_id> and <user_key> with actual values when making the API requests. Also, make sure to handle the responses accordingly in your code.