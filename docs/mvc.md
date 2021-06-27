1) Benefits of using MVC

- Common framework for people to work within
- Makes it easier to know where to look for bugs/add functionality
- Makes testing easier

2) What does MVC stand for and what does each one mean
M = Model
- Responsible for reading/Writing/accessing data
- Common use cases are connecting to DB
- This COULD be connecting other services/APIs

V = View
- Render data/generate the things for the user
- Typically HTML, but can be JSON, text, etc.
- Should have very little logic

C = Controller
- "Air traffic controller"
- Routing request between different pieces of code based on responses it gets
