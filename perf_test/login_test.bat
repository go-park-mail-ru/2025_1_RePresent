@echo off
Remove-Item login_targets.json -Force
go run .\generate_users.go

echo Running Vegeta login test...
if not exist results mkdir results

Get-Content ./login_targets.json | vegeta attack -format=json -rate=100 -duration=60s -connections=10 -max-connections=100 -max-workers=100 -output=results/login_result_100_RPS.bin
vegeta report -type=text results/login_result_100_RPS > results/login_report_100_RPS.txt
vegeta plot results/login_result_100_RPS > results/login_result_plot_100_RPS.html
.\results\login_result_plot_100_RPS.html
