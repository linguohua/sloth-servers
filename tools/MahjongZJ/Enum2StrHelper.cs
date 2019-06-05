using System.Text;
using mahjong;

namespace MahjongTest
{
    public class Enum2StrHelper
    {
        public static string GreatWinType2String(int greatWinGreatWinType)
        {
            var sb = new StringBuilder();
            if (0 != (greatWinGreatWinType & (int) GreatWinType.enumGreatWinType_ChowPongKong))
            {
                sb.Append("独钓,");
            }
            if (0 != (greatWinGreatWinType & (int)GreatWinType.enumGreatWinType_FinalDraw))
            {
                sb.Append("海底捞月,");
            }

            if (0 != (greatWinGreatWinType & (int)GreatWinType.enumGreatWinType_PongKong))
            {
                sb.Append("碰碰胡,");
            }
            if (0 != (greatWinGreatWinType & (int)GreatWinType.enumGreatWinType_PureSame))
            {
                sb.Append("清一色,");
            }
            if (0 != (greatWinGreatWinType & (int)GreatWinType.enumGreatWinType_MixedSame))
            {
                sb.Append("混一色,");
            }
            if (0 != (greatWinGreatWinType & (int)GreatWinType.enumGreatWinType_ClearFront))
            {
                sb.Append("大门清,");
            }
            if (0 != (greatWinGreatWinType & (int)GreatWinType.enumGreatWinType_SevenPair))
            {
                sb.Append("七对,");
            }
            if (0 != (greatWinGreatWinType & (int)GreatWinType.enumGreatWinType_GreatSevenPair))
            {
                sb.Append("豪华七对,");
            }
            if (0 != (greatWinGreatWinType & (int)GreatWinType.enumGreatWinType_Heaven))
            {
                sb.Append("天胡,");
            }
            if (0 != (greatWinGreatWinType & (int)GreatWinType.enumGreatWinType_AfterConcealedKong))
            {
                sb.Append("暗杠胡,");
            }
            if (0 != (greatWinGreatWinType & (int)GreatWinType.enumGreatWinType_AfterExposedKong))
            {
                sb.Append("明杠胡,");
            }
            if (0 != (greatWinGreatWinType & (int)GreatWinType.enumGreatWinType_Richi))
            {
                sb.Append("起手报听,");
            }

            return sb.ToString();
        }

        public static string WinType2String(int playerScoreWinType)
        {
            var result = "";
            switch (playerScoreWinType)
            {
                case (int)HandOverType.enumHandOverType_Win_Chuck:
                    result = "吃铳胡";
                    break;
                case (int)HandOverType.enumHandOverType_Win_SelfDrawn:
                    result = "自摸胡";
                    break;
                case (int)HandOverType.enumHandOverType_None:
                    result = "";
                    break;
                case (int)HandOverType.enumHandOverType_Chucker:
                    result = "放铳";
                    break;
            }
            return result;
        }

        public static string EndType2String(int msgEndType)
        {
            var result = "";
            switch (msgEndType)
            {
                case (int)HandOverType.enumHandOverType_Win_Chuck:
                    result = "放铳胡牌";
                    break;
                case (int)HandOverType.enumHandOverType_Win_SelfDrawn:
                    result = "自摸胡牌";
                    break;
                case (int)HandOverType.enumHandOverType_None:
                    result = "流局";
                    break;
            }
            return result;
        }

        public static  string ChairId2Name(int chairId)
        {
            var result = "E";
            switch (chairId)
            {
                case 0:
                    result = "A";
                    break;
                case 1:
                    result = "B";
                    break;
                case 2:
                    result = "C";
                    break;
                case 3:
                    result = "D";
                    break;
            }
            return result;
        }

        public static string MiniWinType2String(int miniWinMiniWinType)
        {
            var sb = new StringBuilder();
            if (0 != (miniWinMiniWinType & (int)MiniWinType.enumMiniWinType_SelfDraw))
            {
                sb.Append("自摸X2,");
            }

            if (0 != (miniWinMiniWinType & (int)MiniWinType.enumMiniWinType_Continuous_Banker))
            {
                sb.Append("连庄X2,");
            }

            if (0 != (miniWinMiniWinType & (int)MiniWinType.enumMiniWinType_NoFlowers))
            {
                sb.Append("吃椪杠10花,");
            }

            if (0 != (miniWinMiniWinType & (int)MiniWinType.enumMiniWinType_Kong2Discard))
            {
                sb.Append("杠冲X2,");
            }

            if (0 != (miniWinMiniWinType & (int)MiniWinType.enumMiniWinType_Kong2SelfDraw))
            {
                sb.Append("杠开X2,");
            }

            if (0 != (miniWinMiniWinType & (int)MiniWinType.enumMiniWinType_SecondFrontClear))
            {
                sb.Append("小门清X2,");
            }

            return sb.ToString();
        }
    }
}
