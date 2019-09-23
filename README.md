# Test task for Onefootball

## Notes
Part of the task was to scan given API url in order to find ids for given teams.  
It looked like a big overhead to scan your API searching required teams so i spent some 
time investigating your public API and found search query: `https://api.onefootball.com/entity-index-api/v1/archive/team/en-a.json`  
That allowed me to minimize number of required HTTP calls to `2N` where N is the number of given teams

I have intentionally ommited tests as i think there is no reason to test iterations over maps and arrays for such simple task.  

If you don't mind, you can take a look on another test task, that was done by me almost a year ago. It covers the client-server app development topic and writing tests: https://github.com/ivch/testms

## Requirements
- Docker 17.05+

## Build
- `make build`

## Run
- `make run`

## Clean
Cleans docker images created by `make run` or `make build`
- `make clean`
