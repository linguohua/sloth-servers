
namespace MahjongTest
{
    public class DealCfgSimple
    {
        public static  readonly  DealCfgSimple Empty = new DealCfgSimple() {Name=""};

        public string Name { get; set; }

        public int PlayerCount { get; set; }

        public override string ToString()
        {
            return $"{Name},players:{PlayerCount}";
        }
    }
}
