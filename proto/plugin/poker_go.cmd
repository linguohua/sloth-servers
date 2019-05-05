protoc --proto_path=../../proto --go_out=../../gosrc/pokerface  ../../proto/game_pokerface.proto ../../proto/game_pokerface_split2.proto ../../proto/game_pokerface_replay.proto ../../proto/game_pokerface_s2s.proto

@REM 跑得快
protoc --proto_path=../../proto --go_out=../../gosrc/prserver/prunfast ../../proto/game_pokerface_rf.proto

@pause
