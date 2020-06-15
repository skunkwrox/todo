

Data structures

    Person:
        string name
        int64 id
        string email

    Task:
        string title
        string description
        int64 id
        Datetime due_by
        Datetime create
        Datetime last_updated
        Person assignedTo
        Priority priority


REST interface:

    End point                       method      Action
    /person                         GET         returns all persons
                                                  (ids and names)
    /person/<person_id>             GET         return the details for a person
    /person                         POST        add a new person
    /person/<person_id>             PUT         update a person's information
    /person/<person_id>             DELETE      remove a person
    
    
    End point                       method      Action
    /task                           GET         returns all tasks
    /task/<task_id>                 GET         returns the details for a task
    /task                           POST        add a task
    /task/<task_id>                 PUT         updates a task
    /task/<task_id>                 DELETE      deletes a task                                              


Packages:

    go get -u github.com/gorilla/mux
    go get -u github.com/jmoiron/sqlx
    go get -u github.com/lib/pq