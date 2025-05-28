.PHONY: test

test:
	go test -coverprofile=coverage.out ./... | tee full.log && \
	awk '/coverage:/ { \
		for (i=1; i<=NF; i++) { \
			if ($$i ~ /coverage:/) { \
				pct=$$(i+1); \
				gsub(/%/, "", pct); \
				if (pct+0 > 0) { \
					print $$(2); \
				} \
			} \
		} \
	}' full.log > covered_pkgs.txt && \
	echo "Test cover was successful got" && \
	xargs -a covered_pkgs.txt go test -cover -coverprofile=coverage_filtered.out && \
	go tool cover -func=coverage_filtered.out | grep total

cover:
	xargs -a covered_pkgs.txt go test -cover -coverprofile=coverage_filtered.out && \
    go tool cover -func=coverage_filtered.out | grep total