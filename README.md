# TODO
Solution for the [task-tracker](https://roadmap.sh/projects/task-tracker)

##How  to use

Clone the repository and run the following command:

```bash
git clone https://github.com/gulmix/TODO.git
```

Run the following command to build and run the project:

```bash
go build main.go

# To add a task
go run main.go add "Buy groceries"

# To update a task
go run main.go update 1 "Buy groceries and cook dinner"

# To delete a task
go run maing.go delete 1

# To mark a task as in progress/done
go run maing.go mark-in-progress 1
go run maing.go mark-done 1

# To list all tasks
go run maing.go list
go run maing.go list done
go run maing.go list in-progress
```
