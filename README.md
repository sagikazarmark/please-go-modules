# Please Go Modules experiment

Experiments with [Please](https://please.build) and Go modules.


## Usage

```
./pleasew build --rebuild //:bin -v debug --show_all_output
./pleasew run //:bin -v debug --show_all_output
```


**Note:** on the first run it fails (probably because I don't know please very well yet):

```
Build stopped after 1.52s. 1 target failed:
    //:gomod
Attempting to record rule hash: cannot calculate hash for plz-out/gen/pkg/darwin_amd64/github.com/pkg/errors/errors.a: file does not exist
```


## Notes

- This is probably a great example of what should never go into production!
- It used `go mod vendor` before, but it turned out to be a dead end: it requires all the source files in the tree and will only work with Go 1.15
- It uses `go mod download` and outputs the `pkg` directory from GOPATH...don't know if it's a good idea or not
- Currently still requires dependencies to be manually listed
- How automatically creating targets that can be referenced by other targets is still a mystery to me

