@echo off
Remove-Item signup_targets.json -Force
go run .\generate_users.go

echo Running Vegeta signup test...
if not exist results mkdir results

Get-Content ./signup_targets.json | vegeta attack -format=json -rate=100 -duration=60s -connections=10 -max-connections=100 -max-workers=100 -output=results/signup_result_100_RPS.bin
vegeta report -type=text results/signup_result_100_RPS > results/signup_report_100_RPS.txt
vegeta plot results/signup_result_100_RPS > results/signup_result_plot_100_RPS.html
.\results\signup_result_plot_100_RPS.html