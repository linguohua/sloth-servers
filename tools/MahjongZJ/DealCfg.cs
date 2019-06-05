using CsvHelper;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;

namespace MahjongTest
{
    public class DealCfg
    {
        public const int MAX_FLOWER_COUNT = 36;

        public class DealCfgTag
        {
            public DealCfg DealCfg;
            public int Tile;
        }

        public DealCfg(int i, DealCfgWnd owner)
        {
            Owner = owner;
            IsBanker = i == 0;
            Index = i;
        }

        public DealCfgWnd Owner;
        public int Index;
        public bool IsBanker { get; }
        public List<int> TilesHand = new List<int>();

        public Button[] ButtonsHand = new Button[14];
        public Button[] ButtonsFlower = new Button[MAX_FLOWER_COUNT];
        public GroupBox GroupBox;

        public void HideAllButtons()
        {
            foreach (var button in ButtonsHand)
            {
                button.Visibility = Visibility.Collapsed;
            }

            foreach (var button in ButtonsFlower)
            {
                button.Visibility = Visibility.Collapsed;
            }
        }

        public void Tiles2Ui()
        {
            HideAllButtons();

            int i = 0;
            foreach (var tile in TilesHand)
            {
                var btn = ButtonsHand[i];
                btn.Tag = new DealCfgTag() { DealCfg = this, Tile = tile };
                btn.Content = new Image() { Source = Owner._owner.ImageDict[tile] };
                btn.Visibility = Visibility.Visible;
                i++;
            }
            

            var tag = "庄家";
            if (Index > 0)
            {
                tag = "闲家" + Index;
            }

            GroupBox.Header = $"{tag}(手:{TilesHand.Count})";
        }

        public void WriteCsv(CsvWriter csv)
        {
            csv.WriteField(""); // userID

            var handStr = ToHandString();
            csv.WriteField(handStr);
            
            csv.WriteField(""); // 动作提示
        }

        public string ToHandString()
        {
            var hands = TilesHand.ToArray();
            Array.Sort(hands);

            var sb = new StringBuilder();
            foreach (var hand in hands)
            {
                sb.Append(Owner._owner.IdNames[hand]);
                sb.Append(",");
            }
            return sb.ToString();
        }

        public void ReadCsv(CsvReader csvReader)
        {
            var filedBegin = Index * 3 + 2;

            var handTilesStrs = csvReader.GetField(filedBegin + 1);
            var handTilesStrArray = handTilesStrs.Split(',', '，', ' ', '\t');
            var handTotal = 13;
            if (IsBanker)
            {
                handTotal = 14;
            }

            foreach (var s in handTilesStrArray)
            {
                if (!string.IsNullOrWhiteSpace(s) && TilesHand.Count() < handTotal)
                {
                    var tid = Owner._owner.NameIds[s];
                    if (Owner.WallTiles[tid] > 0)
                    {
                        Owner.WallTiles[tid]--;
                        TilesHand.Add(tid);
                    }
                }
            }
        }
    }
}
