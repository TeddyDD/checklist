# Checklist

Checklists are important aid for pilots, why not use them in IT?

## Example

Checklist for updating void packages

```
- noone sent PR for this package yet
- read changelog
- update checksums
- checksums is correct `./xbps_src extract $1`
- xlint returns no error `xlint srcpkgs/$1/template`
- package is building
- commiting on correct branch `[ "$(git rev-parse --abbrev-ref HEAD)" = "$1" ]`
```

```sh
ckl void lf && xbump
```

