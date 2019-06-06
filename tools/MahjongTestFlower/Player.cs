using System;
using System.Windows;
using mahjong;
using WebSocketSharp;

namespace MahjongTest
{
    public class Player : IDisposable
    {
        
        public WebSocket Ws { get; }
        public string Name { get; }

        public string UserId { get; }
        public TileStackWnd MyWnd { get; }

        public MainWindow MWnd { get; }
        public Player(string name, string userId, string roomNumber, TileStackWnd myWnd, MainWindow mWnd)
        {
            Name = name;
            MyWnd = myWnd;
            UserId = userId;
            var url = $"{ProgramConfig.ServerUrl}/ws/monkey?userID={userId}&roomNumber={roomNumber}";
            if (url.StartsWith("https://"))
            {
                url = url.Replace("https", "wss");
            }
            else
            {
                url = url.Replace("http", "ws");
            }
 
            Ws = new WebSocket(string.Format(url, userId, roomNumber));
            MyWnd.SetPlayer(this);
            MWnd = mWnd;
        }

        public int ChairId
        {
            get;
            set;
        }

        public void Connect()
        {
            Ws.OnMessage += OnMessageThread;
            Ws.OnClose += OnCloseThread;

            Ws.Connect();

            MyWnd.Reset2New();
        }

        private void OnCloseThread(object sender, CloseEventArgs e)
        {
            Dispose();

            Action a = OnPlayerDisconnected;

            MyWnd.Dispatcher.Invoke(a);
        }

        private void OnPlayerDisconnected()
        {
            MyWnd.Visibility = Visibility.Hidden;
            MWnd.OnPlayerDisconnected(this);
        }

        private void OnMessageThread(object sender, MessageEventArgs messageEventArgs)
        {
            var gmsg = messageEventArgs.RawData.ToProto<GameMessage>();
            //Console.WriteLine($"player got message, ops:{gmsg.Ops}");

            Action a = () =>
            {
                OnServerMessage(this, gmsg);
            };
            MyWnd.Dispatcher.Invoke(a);
        }

        private static void OnServerMessage(Player player, GameMessage gmsg)
        {
            switch (gmsg.Ops)
            {
                case (int)MessageCode.OPActionAllowed:
                    {
                        var msg = gmsg.Data.ToProto<MsgAllowPlayerAction>();
                        OnServerMessageActionAllowed(player, msg);
                    }
                    break;
                case (int)MessageCode.OPReActionAllowed:
                    {
                        var msg = gmsg.Data.ToProto<MsgAllowPlayerReAction>();
                        OnServerMessageReActionAllowed(player, msg);
                    }
                    break;
                case (int)MessageCode.OPActionResultNotify:
                    {
                        var msg = gmsg.Data.ToProto<MsgActionResultNotify>();
                        OnServerMessageActionResultNotify(player, msg);
                    }
                    break;
                case (int)MessageCode.OPDeal:
                    {
                        var msg = gmsg.Data.ToProto<MsgDeal>();
                        OnServerMessageDeal(player, msg);
                    }
                    break;
                case (int)MessageCode.OPHandOver:
                    {
                        var msg = gmsg.Data.ToProto<MsgHandOver>();
                        OnServerMessageHandScore(player, msg);
                    }
                    break;
                case (int)MessageCode.OPPlayerEnterRoom:
                    {
                        var msg = gmsg.Data.ToProto<MsgEnterRoomResult>();
                        OnServerMessageEnterRoom(player, msg);
                    }
                    break;
                case (int)MessageCode.OPRoomUpdate:
                    {
                        var msg = gmsg.Data.ToProto<MsgRoomInfo>();
                        OnServerMessageRoomUpdate(player, msg);
                    }
                    break;
                case (int)MessageCode.OPRoomShowTips:
                    {
                        var msg = gmsg.Data.ToProto<MsgRoomShowTips>();
                        OnServerMessageRoomShowTips(player, msg);
                    }
                    break;
                case (int)MessageCode.OPDisbandNotify:
                    {
                        var msg = gmsg.Data.ToProto<MsgDisbandNotify>();
                        OnServerDisbandNotify(player, msg);
                    }
                    break;
            }
        }

        private static void OnServerDisbandNotify(Player player, MsgDisbandNotify msg)
        {
            player.MyWnd.OnDisbandNotify(msg);
        }

        private static void OnServerMessageRoomShowTips(Player player, MsgRoomShowTips msg)
        {
            // 获得服务器分配的chair id
            player.MyWnd.OnShowRoomTips(msg);
        }

        private static void OnServerMessageRoomUpdate(Player player, MsgRoomInfo msg)
        {
            // 获得服务器分配的chair id
            foreach (var playerInfo in msg.players)
            {
                if (playerInfo.userID == player.UserId)
                {
                    player.ChairId = playerInfo.chairID;
                }
            }
        }

        private static void OnServerMessageEnterRoom(Player player, MsgEnterRoomResult msg)
        {
            player.MyWnd.OnEnterRoom(msg);
        }

        private static void OnServerMessageHandScore(Player player, MsgHandOver msg)
        {
            player.MyWnd.OnHandScore(msg);
        }

        private static void OnServerMessageActionResultNotify(Player player, MsgActionResultNotify msg)
        {
            //throw new NotImplementedException();
            if (msg.targetChairID == player.ChairId)
            {
                // my result
                player.MyWnd.OnActionResult(msg);
            }
            else
            {
                if (msg.action != (int)ActionType.enumActionType_FirstReadyHand)
                    player.MyWnd.CancelAllowedAction();
            }
        }

        private static void OnServerMessageActionAllowed(Player player, MsgAllowPlayerAction msg)
        {
            //throw new NotImplementedException();
            if (msg.actionChairID == player.ChairId)
            {
                // my actions
                player.MyWnd.OnAllowedActions(msg);
            }
        }
        private static void OnServerMessageReActionAllowed(Player player, MsgAllowPlayerReAction msg)
        {
            //throw new NotImplementedException();
            if (msg.actionChairID == player.ChairId)
            {
                // my actions
                player.MyWnd.OnAllowedReActions(msg);
            }
        }

        private static void OnServerMessageDeal(Player player, MsgDeal msg)
        {
            //throw new NotImplementedException();
            player.MyWnd.ResetPlayStatus();
            player.MyWnd.OnDeal(msg);
        }

        public void Dispose()
        {
            ((IDisposable)Ws)?.Dispose();
        }

        public void SendMessage(int opAction, byte[] toBytes)
        {
            var gmsg = new GameMessage
            {
                Ops = (int)opAction,
                Data = toBytes
            };
            var msgBytes = gmsg.ToBytes();
            Ws?.Send(msgBytes);
        }

        public void SendReady2Server()
        {
            MyWnd.SendReady2Server();
        }
    }
}
