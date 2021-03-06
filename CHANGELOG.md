## Version 1.2.1 (Release date: 2017-04-04)
([1212d45](https://github.com/nytimes/video-transcoding-api/commit/1212d45)) provider/encodingcom: fix bug in newly introduced adjustSize function 




## Version 1.2.0 (Release date: 2017-04-04)
([9cbe954](https://github.com/nytimes/video-transcoding-api/commit/9cbe954)) Fix gometalinter violations 

([5776807](https://github.com/nytimes/video-transcoding-api/commit/5776807)) provider/encodingcom: provide the proper information on JobStatus 


([41ae17b](https://github.com/nytimes/video-transcoding-api/commit/41ae17b)) new fields for README 

([05dc0ea](https://github.com/nytimes/video-transcoding-api/commit/05dc0ea)) adding flexibility for s3 storage region and the ability to encode in any region of aws or gcp 


([6f5cf95](https://github.com/nytimes/video-transcoding-api/commit/6f5cf95)) travis: use .x notation for Go 1.8 


([8eaa73e](https://github.com/nytimes/video-transcoding-api/commit/8eaa73e)) code cleanup 

([6deee07](https://github.com/nytimes/video-transcoding-api/commit/6deee07)) removing debugging code 

([1d8e517](https://github.com/nytimes/video-transcoding-api/commit/1d8e517)) linter caught some dead code and other things, fixing this 

([e93afc3](https://github.com/nytimes/video-transcoding-api/commit/e93afc3)) code cleanup 

([b12c739](https://github.com/nytimes/video-transcoding-api/commit/b12c739)) more test coverage and code cleanup 

([0c1c815](https://github.com/nytimes/video-transcoding-api/commit/0c1c815)) adding test code 

([da1a24e](https://github.com/nytimes/video-transcoding-api/commit/da1a24e)) travis: update go 1.8 to rc3 

([1b3bd2d](https://github.com/nytimes/video-transcoding-api/commit/1b3bd2d)) adding Bitmovin to readme 

([ccb4321](https://github.com/nytimes/video-transcoding-api/commit/ccb4321)) jobs submit but does not respond with proper json, debugging code in here for now, manifest gen works for hls 

([4af9b33](https://github.com/nytimes/video-transcoding-api/commit/4af9b33)) changing import directory 

([7a0a658](https://github.com/nytimes/video-transcoding-api/commit/7a0a658)) need to register bitmovin 

([2900b79](https://github.com/nytimes/video-transcoding-api/commit/2900b79)) fleshing out rest of methods 

([ec748c5](https://github.com/nytimes/video-transcoding-api/commit/ec748c5)) interstitial commit 

([29272e9](https://github.com/nytimes/video-transcoding-api/commit/29272e9)) change in logic, adding a lot of steps 

([e2c10d9](https://github.com/nytimes/video-transcoding-api/commit/e2c10d9)) adding environment variables for bitmovin aws credentials 

([2fbc577](https://github.com/nytimes/video-transcoding-api/commit/2fbc577)) s3 url parsing 

([3a05910](https://github.com/nytimes/video-transcoding-api/commit/3a05910)) interstitial commit 

([50b89f3](https://github.com/nytimes/video-transcoding-api/commit/50b89f3)) interstitial commit 

([595dce8](https://github.com/nytimes/video-transcoding-api/commit/595dce8)) implementing the interface 

([a34faef](https://github.com/nytimes/video-transcoding-api/commit/a34faef)) adding rest of config for bitmovin clinet 

([02e97d3](https://github.com/nytimes/video-transcoding-api/commit/02e97d3)) initial config changes for bitmovin 

([003fcd1](https://github.com/nytimes/video-transcoding-api/commit/003fcd1)) adding stubs for functions, capapbilities and first test 

([6f424e6](https://github.com/nytimes/video-transcoding-api/commit/6f424e6)) initial commit 





## Version 1.1.2 (Release date: 2017-01-25)

([f32e347](https://github.com/nytimes/video-transcoding-api/commit/f32e347)) Changes the Zencoder wrapper to use the job status found in the jobDetails for the JobState 

([58b5e83](https://github.com/nytimes/video-transcoding-api/commit/58b5e83)) travis: update go 1.8 to rc2 

([9357eaa](https://github.com/nytimes/video-transcoding-api/commit/9357eaa)) db/redis: update go-redis 

([17730c0](https://github.com/nytimes/video-transcoding-api/commit/17730c0)) drone: run integration tests after deploying 

([fa7488e](https://github.com/nytimes/video-transcoding-api/commit/fa7488e)) travis: update Go 1.8 to rc1 

([ba8e873](https://github.com/nytimes/video-transcoding-api/commit/ba8e873)) readme: don't use an actual IP address in example 

([43e243e](https://github.com/nytimes/video-transcoding-api/commit/43e243e)) Update swagger.json 

([f4bd0eb](https://github.com/nytimes/video-transcoding-api/commit/f4bd0eb)) doc: include zencoder 




## Version 1.1.1 (Release date: 2017-01-06)

([98ba09f](https://github.com/nytimes/video-transcoding-api/commit/98ba09f)) travis: use .x syntax to ensure latest 1.7 

([24d9528](https://github.com/nytimes/video-transcoding-api/commit/24d9528)) db/redis/storage: support float64 




## Version 1.1.0 (Release date: 2016-12-22)

([ec42742](https://github.com/nytimes/video-transcoding-api/commit/ec42742)) encodingcom: Returns converted file size 

([12b0a6d](https://github.com/nytimes/video-transcoding-api/commit/12b0a6d)) Update gops 

([0c1152c](https://github.com/nytimes/video-transcoding-api/commit/0c1152c)) encodingcom: Returns converted file size 





## Version 1.0.9 (Release date: 2016-12-19)



## Version 1.0.8-rc (Release date: 2016-12-19)

([071a2d8](https://github.com/nytimes/video-transcoding-api/commit/071a2d8)) preset: avoid creating PresetMap when ProviderMapping is empty 

([314930b](https://github.com/nytimes/video-transcoding-api/commit/314930b)) preset: improve if statement 




## Version 1.0.7 (Release date: 2016-12-16)



## Version 1.0.6-rc (Release date: 2016-12-16)

([d863498](https://github.com/nytimes//commit/d863498)) Add filesize rendition info to Zencoder 


([6de4a74](https://github.com/nytimes//commit/6de4a74)) preset: bubble up the error when creating a preset 

([39c0992](https://github.com/nytimes//commit/39c0992)) service/preset: remove logging and fix comments 

([90d24dc](https://github.com/nytimes//commit/90d24dc)) service/presetmap: create or update existent presetmap when creating new presets 




## Version 1.0.5 (Release date: 2016-12-08)
([c53ea3d](https://github.com/nytimes/video-transcoding-api/commit/c53ea3d)) provider/elementalconductor: fix compatibility with encoding-wrapper 

([6e951cc](https://github.com/nytimes/video-transcoding-api/commit/6e951cc)) travis: add Go 1.8beta1 

([0295924](https://github.com/nytimes/video-transcoding-api/commit/0295924)) swagger: add test to increase package coverage 


([c60a8f3](https://github.com/nytimes/video-transcoding-api/commit/c60a8f3)) zencoder: use constants for Job State 

([44e3bd5](https://github.com/nytimes/video-transcoding-api/commit/44e3bd5)) zencoder: set progress to 100 when job status is finished (fixes #170) 




## Version 1.0.4 (Release date: 2016-12-05)

([b220d80](https://github.com/nytimes/video-transcoding-api/commit/b220d80)) zencoder: bugfix on duration being reported 


([4566b29](https://github.com/nytimes/video-transcoding-api/commit/4566b29)) Revert "encodingcom: remove dead code when creating a encoding.com Format" 




## Version 1.0.3 (Release date: 2016-12-02)

([db75994](https://github.com/nytimes/video-transcoding-api/commit/db75994)) Protect against possible stray colon in error message 

([a56d8da](https://github.com/nytimes/video-transcoding-api/commit/a56d8da)) Revert "Populate status message" 


([4df203f](https://github.com/nytimes/video-transcoding-api/commit/4df203f)) Populate job status with detailed status message 




## Version 1.0.2-rc (Release date: 2016-12-01)
([98d5b5a](https://github.com/nytimes/video-transcoding-api/commit/98d5b5a)) Makefile: use CI_TAG in `make live` when available 

([d0c5e27](https://github.com/nytimes/video-transcoding-api/commit/d0c5e27)) travis: update Go 




## Version 1.0.1-rc (Release date: 2016-12-01)

([a6112be](https://github.com/nytimes/video-transcoding-api/commit/a6112be)) zencoder: consider finished outputs with no format and m3u8 suffix as m3u8 container (refs #161) 


([d68b5e2](https://github.com/nytimes/video-transcoding-api/commit/d68b5e2)) provider/zencoder: use GetVodUsage in Healthcheck 




([58bd6e8](https://github.com/nytimes/video-transcoding-api/commit/58bd6e8)) makefile: fix variable replacement for makefile scheme 

([766e434](https://github.com/nytimes/video-transcoding-api/commit/766e434)) build: detach stg and prod deployment by using 'rc' on tag name 

([b3d1b63](https://github.com/nytimes/video-transcoding-api/commit/b3d1b63)) encodingcom: remove dead code when creating a encoding.com Format 




## Version 1.0.0 (Release date: 2016-11-23)



## Version 0.1.6 (Release date: 2016-11-23)
([dbf6d29](https://github.com/nytimes/video-transcoding-api/commit/dbf6d29)) zencoder: fix hls path (close #157) 




## Version 0.1.5 (Release date: 2016-11-23)
([bf7521b](https://github.com/nytimes/video-transcoding-api/commit/bf7521b)) zencoder: avoid concatenating 'hls' path to hls output (refs #157) 





## Version 0.1.4 (Release date: 2016-11-23)




## Version 0.1.2 (Release date: 2016-11-21)

([f3e2435](https://github.com/nytimes/video-transcoding-api/commit/f3e2435)) zencoder: add PrepareForSegmenting: 'hls' for mp4's that matches with HLS 

([e967ab0](https://github.com/nytimes/video-transcoding-api/commit/e967ab0)) zencoder: fix golint complain 

([2a0a755](https://github.com/nytimes/video-transcoding-api/commit/2a0a755)) zencoder: bugfix on isOutputCompatible() method 

([2d69d00](https://github.com/nytimes/video-transcoding-api/commit/2d69d00)) zencoder: raise errors gracefully 

([59e6e75](https://github.com/nytimes/video-transcoding-api/commit/59e6e75)) zencoder: reuse mp4 outputs for transmuxing hls outputs (close #151) 


([7f0d41c](https://github.com/nytimes/video-transcoding-api/commit/7f0d41c)) Send logging and error reporting via agent 




## Version 0.1.1 (Release date: 2016-11-16)

([225d9ff](https://github.com/nytimes/video-transcoding-api/commit/225d9ff)) zencoder: make all zencoder uploads public 


([37acd1c](https://github.com/nytimes/video-transcoding-api/commit/37acd1c)) zencoder: normalize hls output names based on encoding.com implementation 


([8666837](https://github.com/nytimes/video-transcoding-api/commit/8666837)) db: remove unused struct from stub_test 

([d6a1595](https://github.com/nytimes/video-transcoding-api/commit/d6a1595)) db/redis: add structs to stub_test to avoid db dependency on redis storage 

([a70a602](https://github.com/nytimes/video-transcoding-api/commit/a70a602)) db/redis: add more tests for FieldMap() method 

([4231353](https://github.com/nytimes/video-transcoding-api/commit/4231353)) db/redis: add test for FieldMap() method 