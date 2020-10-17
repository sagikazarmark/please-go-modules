# Please Go modules example

Regenerate rules:

```bash
plz run --in_wd //cmd/godeps -- -dir third_party/go -base example -clean -subinclude "//build_defs"
```
