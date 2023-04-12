# Replace

It works locally, refering the package arent published yeat

## Example:
> go mod edit -replace github.com/mmfalcao/go-expert/packaging/3/math=../math
Then:
> go mod tidy
> go run *.go