run:
	go run main.go >> output.txt
merge-csv:
	go run cli-apps/csv-merger/app.go