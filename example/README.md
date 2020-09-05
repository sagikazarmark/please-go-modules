# Please Go modules example

Regenerate rules:

```bash
plz build //cmd/gogetgen
../plz-out/bin/cmd/gogetgen/gogetgen -dir example/third_party/go -clean -genpkg -subinclude "//build_defs"
rm -rf third_party
mv example/third_party .
rm -rf example
```
