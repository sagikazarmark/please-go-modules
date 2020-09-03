# Please Go modules example

Regenerate rules:

```bash
plz run --in_wd //cmd/gogetgen -- -dir example/third_party/go -clean -genpkg
rm -rf third_party
mv example/third_party .
rm -rf example
```
