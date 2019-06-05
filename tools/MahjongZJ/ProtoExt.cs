using System;
using System.IO;
using mahjong;
using ProtoBuf;

namespace MahjongTest
{
    public static class ProtoExt
    {
        public static byte[] ToBytes<T>(this T proto)
        {
            if (proto == null)
                return null;

            using (var ms = new MemoryStream())
            {
                Serializer.Serialize(ms, proto);
                return ms.ToArray();
            }
        }

        public static GameMessage ToMessage<T>(this T proto, int ops)
        {
            return ToMessage(proto, ops, 0, 0);
        }

        public static GameMessage ToMessage<T>(this T proto, int ops, long playerId)
        {
            return ToMessage(proto, ops, 0, playerId);
        }

        public static GameMessage ToMessage<T>(this T proto, int ops, int serverid)
        {
            return ToMessage(proto, ops, serverid, 0);
        }

        public static GameMessage ToMessage<T>(this T proto, int ops, int serverId, long playerId)
        {
            var ret = new GameMessage
            {
                Ops = ops,
                Data = proto.ToBytes()
            };


            return ret;
        }

        public static T ToProto<T>(this Stream stream)
        {
            if (stream == null) return default(T);
            return Serializer.Deserialize<T>(stream);
        }

        public static T ToProto<T>(this byte[] data)
        {
            if (data == null || data.Length == 0) return default(T);
            try
            {
                using (var ms = new MemoryStream(data))
                {
                    return Serializer.Deserialize<T>(ms);
                }
            }
            catch (Exception e)
            {
                Console.WriteLine(e.Message);
                return default(T);
            }
        }
    }
}
