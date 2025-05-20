@echo off
Remove-Item signup_targets.json -Force
go run .\generate_users.go

echo Running Vegeta signup test...
if not exist results mkdir results

Get-Content ./signup_targets.json | vegeta attack -format=json -rate=1000 -duration=60s -connections=100 -max-connections=1000 -max-workers=1000 -output=results/signup_result_1k_RPS.bin
vegeta report -type=text results/signup_result_1k_RPS > results/signup_report_1k_RPS.txt
vegeta plot results/signup_result_1k_RPS > results/signup_result_plot_1k_RPS.html
.\results\signup_result_plot_1k_RPS.html