using System;
using System.IO;

namespace MahjongTest
{
    public class ManualProtoCoder
    {
        public const int WireTypeVarint = 0;
        public const int LengthDelimited = 2;
        /// <summary>
        /// 按照protocol buffer协议编码int32
        /// </summary>
        /// <param name="bw">二进制写入器</param>
        /// <param name="val">需要编码的值</param>
        /// <param name="fieldNumber">该值位于Message中的序号</param>
        public static void EncodeInt32(BinaryWriter bw, UInt32 val, int fieldNumber)
        {
            // 先写(field_number << 3) | wire_type
            int tagVariant = (fieldNumber << 3) | WireTypeVarint;
            EncodeVariantUInt32(bw, (UInt32)tagVariant);
            EncodeVariantUInt32(bw, val);
        }
        /// <summary>
        /// 按照protocol buffer协议解码int32
        /// </summary>
        /// <param name="br">二进制读取器</param>
        /// <param name="fieldNumber">该值位于Message中的序号</param>
        /// <param name="val">如果读取成功保存读取到的值</param>
        /// <returns>成功返回true，失败返回false</returns>
        public static bool DecodeInt32(BinaryReader br, int fieldNumber, out UInt32 val)
        {
            val = 0;
            var fw = DecodeVariant(br);
            var fieldNumberX = fw >> 3;
            if (fieldNumberX != fieldNumber)
            {
                return false;
            }

            var wireType = fw & 0x07;
            if (wireType != WireTypeVarint)
            {
                return false;
            }

            val = DecodeVariant(br);

            return true;
        }
        /// <summary>
        /// 按照protocol buffer协议编码字节流
        /// </summary>
        /// <param name="bw">二进制写入器</param>
        /// <param name="data">需要编码的值</param>
        /// <param name="fieldNumber">该值位于Message中的序号</param>
        public static void EncodeBytes(BinaryWriter bw, byte[] data, int fieldNumber)
        {
            // 先写(field_number << 3) | wire_type
            int tagVariant = (fieldNumber << 3) | LengthDelimited;
            EncodeVariantUInt32(bw, (UInt32)tagVariant);
            // 数组长度
            EncodeVariantUInt32(bw, (UInt32)data.Length);
            bw.Write(data);
        }
        /// <summary>
        /// 按照protocol buffer协议解码字节流
        /// </summary>
        /// <param name="br">二进制读取器</param>
        /// <param name="fieldNumber">该值位于Message中的序号</param>
        /// <param name="data">如果读取成功保存读取到的值</param>
        /// <returns>成功返回true，失败返回false</returns>
        public static bool DecodeBytes(BinaryReader br, int fieldNumber, out byte[] data)
        {
            data = null;

            var fw = DecodeVariant(br);
            var fieldNumberX = fw >> 3;
            if (fieldNumberX != fieldNumber)
            {
                return false;
            }

            var wireType = fw & 0x07;
            if (wireType != LengthDelimited)
            {
                return false;
            }

            var length = DecodeVariant(br);
            data = br.ReadBytes((int)length);

            return true;
        }

        private static void EncodeVariantUInt32(BinaryWriter bw, UInt32 variant)
        {
            while (true)
            {
                if (variant > 127)
                {
                    var bytex = 0x80 | (variant & 0x7f);
                    bw.Write((byte)bytex);
                    variant = variant >> 7;
                }
                else
                {
                    bw.Write((byte)variant);
                    break;
                }
            }
        }

        private static UInt32 DecodeVariant(BinaryReader br)
        {
            int n = 0;
            UInt32 v = 0;

            while (true)
            {
                var b = br.ReadByte();
                if (0 == (b & 0x80))
                {
                    // 序列终止
                    UInt32 x = b;
                    v = v | (x << n);
                    break;
                }
                else
                {
                    UInt32 x = (UInt32)(b & 0x7f);
                    v = v | (x << n);
                }
                n += 7;
            }

            return v;
        }

    }
}
