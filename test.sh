# exports environment variables from .env.yaml
# sed 's/: /=/g' replaces : with =
# cloud function formatting requires ':' chars, while bash syntax requires '=' ie $ export ENV_VAR='xxx'
# grep ^# skips lines that are commented out with '#' chars
export $(grep -v '^#' .env.yaml | sed 's/: /=/g' | xargs)

go mod tidy
go build
go test