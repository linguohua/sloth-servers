using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.RegularExpressions;

namespace MahjongTest
{
    internal class Log2X
    {
        internal class Deal
        {
            public string HandTiles { get; set; }

            public string FlowerTiles { get; set; } 

            public string Name { get; set; }

            public int HandTileCount
            {


                get
                {
                    var xx = HandTiles.Split(',');
                    return xx.Count(x => !string.IsNullOrWhiteSpace(x));
                }
            }
        }

        public static HashSet<string> ActionDict { get; } = new HashSet<string>();

        public List<string> Draws { get; } = new List<string>();
        public List<string> ActionLines { get; } = new List<string>();

        public List<Deal> Deals { get;  } = new List<Deal>();

        public string Wind { get; set; }

        public string Banker { get; set; }

        public bool Parse(string log)
        {
            var lines = new List<string>();
            using (StringReader reader = new StringReader(log))
            {

                string rline;
                while ((rline = reader.ReadLine()) != null)
                {
                    // Do something with the line
                    lines.Add(rline);
                }
            }

            if (lines.Count < 1)
                return false;

            var line = lines[0];
            if (line != "[begin]")
                return false;

            var dealLines = lines.Select(x => x).Where(x => x.StartsWith("[deal]") || x.StartsWith("\t(flower)"))
                .ToArray();
            // [deal](A)(hand):1万,2万,3万,1筒,2筒,3筒,8筒,8筒,8筒,1条,2条,3条,7条,
            //  (flower):
            if (dealLines.Length % 2 != 0)
                return false;

            string pattern;
            Regex rgx;
            MatchCollection matches;
            Match match;
            for (var j = 0; j < dealLines.Length-1; j +=2)
            {
                var xhand = dealLines[j];
                var xflower = dealLines[j + 1];

                pattern = @"^\[.+\]\((?<name>.+)\)\(.+\):(?<tiles>.+)$";
                rgx = new Regex(pattern, RegexOptions.IgnoreCase);
                matches = rgx.Matches(xhand);
                if (matches.Count != 1)
                    return false;
                match = matches[0];
                var name = match.Groups["name"].Value;
                var tiles = match.Groups["tiles"].Value;

                if (string.IsNullOrWhiteSpace(name)
                    || string.IsNullOrWhiteSpace(tiles))
                {
                    return false;
                }

                pattern = @"^\s*\(.+\):(?<flowers>.*)$";
                rgx = new Regex(pattern, RegexOptions.IgnoreCase);
                matches = rgx.Matches(xflower);
                if (matches.Count != 1)
                    return false;
                match = matches[0];
                var flowers = match.Groups["flowers"].Value;
                if (string.IsNullOrWhiteSpace(flowers))
                    flowers = "";

                var deal = new Deal() {Name = name, HandTiles = tiles, FlowerTiles = flowers};

                Deals.Add(deal);
            }

            var drawLines = lines.Select(x => x).Where(x => x.StartsWith("[draw]"));
            foreach (var drawLine in drawLines)
            {
                //[draw](A):南
                pattern = @".+:(?<tname>.*)$";
                rgx = new Regex(pattern, RegexOptions.IgnoreCase);
                matches = rgx.Matches(drawLine);
                if (matches.Count != 1)
                    return false;

                match = matches[0];
                var tileName = match.Groups["tname"].Value;
                Draws.Add(tileName);
            }

            ActionLines.AddRange(lines.Select(x => x).Where(IsActionLine));

            var banker = lines.Find(x => x.StartsWith("[bank]"));
            var wind = lines.Find(x => x.StartsWith("[wind]"));
            if (!string.IsNullOrWhiteSpace(banker))
            {
                var begin = banker.IndexOf(":", StringComparison.Ordinal);
                var bname = banker.Substring(begin+1);
                Banker = bname;
            }
            if (!string.IsNullOrWhiteSpace(wind))
            {
                var begin = wind.IndexOf(":", StringComparison.Ordinal);
                var bname = wind.Substring(begin + 1);
                Wind = bname;
            }

            Deals.Sort((x,y)=>x.Name[0] - y.Name[0]);
            var i = 0;
            for (; i < Deals.Count; ++i)
            {
                if (Deals[i].HandTileCount == 14)
                    break;
            }
            var newList = new List<Deal>();
            var length = Deals.Count;
            for (var k = 0; k < length; ++k)
            {
                newList.Add(Deals[(i+k)%length]);
            }
            Deals.Clear();
            Deals.AddRange(newList);
            return true;
        }

        private bool IsActionLine(string s)
        {
            if (ActionDict.Count < 1)
            {
                InitDict();
            }
            var begin = s.IndexOf("[", StringComparison.Ordinal);
            var end = s.IndexOf("]", StringComparison.Ordinal);
            if (begin < 0 || end <= begin)
                return false;

            var actionname = s.Substring(begin+1, end-begin-1);

            return (ActionDict.Contains(actionname));
        }

        private void InitDict()
        { 
            ActionDict.Add("discard");
            ActionDict.Add("chow");
            ActionDict.Add("pong");
            ActionDict.Add("kongExposed");
            ActionDict.Add("kongConcealed");
            ActionDict.Add("triplet2kong");
            ActionDict.Add("richi");
            ActionDict.Add("winchuck");
            ActionDict.Add("skip");
            ActionDict.Add("winselfdraw");
        }
    }
}
