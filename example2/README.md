# Please Go modules example

Regenerate rules:

```bash
cat BUILD.tpl > BUILD
plz run --in_wd //cmd/godeps -- -stdout >> BUILD
```
