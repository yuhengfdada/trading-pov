# Run the Tests
1. Clone the repo & Install Go
2. Under the project root, run `go mod tidy`
3. Run `go test -v ./... > logs`
4. See the test output in `logs`. It should be similar to `sample_log`.

# Create Your Own Tests
## Sample
```go
func TestBehindMin(t *testing.T) {
	lines := setup(t, "behindmin.csv")

	order := makeFIXMsg("1", "1000", "10")
	engine.Order(order)

	engine.SetVolume(5000)
	engine.SetOrderFilledQuantity(100)

	sendEvents(lines)
}
```
1. Put your own dataset in the `datasets` folder.
2. Create a function like the sample inside `app/integrated_test.go`.
3. Change `behindmin.csv` to your own dataset.
4. Create your own client order by modifying the parameters to `makeFIXMsg`.
5. You can also set the engine's state before sending events to it (if you want).
6. Inside the `app` directory, run `go test -run ^your_test_function$ > logs`.
7. See the test output in `app/logs`.