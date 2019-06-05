using CsvHelper;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;

namespace PokerTest
{
    public class DrawCfg
    {
        private readonly DealCfgWnd _owner;
        //public const int MaxDrawCount = 120;

        public DrawCfg(DealCfgWnd owner)
        {
            _owner = owner;
        }

        public List<int> Tiles = new List<int>();
        public List<Button> ButtonsHand = new List<Button>();
        public GroupBox GroupBox;

        public string ToTilesString()
        {
            //var hands = Tiles.ToArray();
            //Array.Sort(hands);

            var sb = new StringBuilder();
            foreach (var tile in Tiles)
            {
                sb.Append(_owner._owner.IdNames[tile]);
                sb.Append(",");
            }
            return sb.ToString();
        }

        public void HideAllButtons()
        {
            foreach (var button in ButtonsHand)
            {
                button.Visibility = Visibility.Collapsed;
            }
        }

        public void Tiles2Ui()
        {
            HideAllButtons();

            int i = 0;
            foreach (var tile in Tiles)
            {
                Button btn = null;
                if (i < ButtonsHand.Count)
                {
                    btn = ButtonsHand[i];
                }
                else
                {
                    btn = _owner.AddDrawUIButton();
                    ButtonsHand.Add(btn);
                }

                btn.Tag = tile;
                btn.Content = new Image() { Source = _owner._owner.ImageDict[tile] };
                btn.Visibility = Visibility.Visible;
                i++;
            }


            GroupBox.Header = $"抽牌序列:{Tiles.Count}";
        }

        public void ReadCsv(CsvReader csvReader)
        {
            var drawSeqStrs = csvReader.GetField(14);
            var drawSeqStrArray = drawSeqStrs.Split(',', '，', ' ', '\t');

            foreach (var s in drawSeqStrArray)
            {
                if (!string.IsNullOrWhiteSpace(s))
                {
                    var tid = _owner._owner.NameIds[s];
                    if (_owner._wallTiles[tid] > 0)
                    {
                        _owner._wallTiles[tid]--;
                        Tiles.Add(tid);
                    }

                }
            }
        }
    }

}
