# Workspace

Create an workspace environment

## Example:
> go work init ./math ./system
Then will generate a file:
> go.work
So will be abble to run your code:
> go run system/*.go

## Release at 1.18

## Caution
When you use the command "go mod tidy" should expect an error

How do you fix this
- First Approach:
  > go get moduleFromWeb
- Second Approach: Publish your module into your github
- Third Approach: this will ignore modules that Go couldn't not fix
  > go mod tidy -e