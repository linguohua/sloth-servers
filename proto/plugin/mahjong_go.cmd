protoc --proto_path=../../proto --go_out=../../gosrc/mahjong  ../../proto/game_mahjong.proto ../../proto/game_mahjong_split2.proto ../../proto/game_mahjong_replay.proto ../../proto/game_mahjong_s2s.proto

@REM 大丰麻将
protoc --proto_path=../../proto --go_out=../../gosrc/dfmjserver/dfmahjong ../../proto/game_mahjong_df.proto

@pause
