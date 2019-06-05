using System;
using System.Collections.Generic;
using pokerface;

namespace PokerTest
{
    /// <summary>
    /// 大丰关张的扑克牌相关计算代码
    /// </summary>
    internal class AgariIndex
    {
        // 保存所有大丰关张中合法的牌型
        private static Dictionary<long, int> agariTable = new Dictionary<long, int>();

        // slots用于把牌按照rank归类，一幅扑克牌从2到JQK到ACE，一共有13个rank，加上大小王一共14个rank
        private static int[] slots = new int[14];

        /// <summary>
        ///  计算手牌的key
        /// 根据牌型计算key，例如JJQQ，得到22，表示两个对子。
        /// </summary>
        /// <param name="hai">手牌列表，例如一个顺子之类</param>
        /// <returns>返回一个int64，如果是lua则返回一个int（lua的int表数范围是48位）</returns>
        static long calcKey(int[] hai)
        {
            for (var i = 0; i < slots.Length; i++)
            {
                slots[i] = 0;
            }

            for (var i = 0; i < hai.Length; i++)
            {
                var h = hai[i];
                slots[h / 4]++;
            }

            for (var i = 0; i < slots.Length; i++)
            {
                if (slots[i] > 4)
                {
                    throw new System.Exception("card type great than 4,card:"+i);
                }
            }

            // 排序，由小到大升序
            Array.Sort(slots);
            long key = 0;
            for (var i = slots.Length-1; i >= 0; i--)
            {
                if (slots[i] == 0)
                {
                    break;
                }

                key = key * 10 + (slots[i]);
            }

            for (var i = 0; i < slots.Length; i++)
            {
                slots[i] = 0;
            }

            for (var i = 0; i < hai.Length; i++)
            {
                var h = hai[i];
                slots[h / 4]++;
            }

            return key;
        }

        /// <summary>
        /// 根据牌列表，构造MsgCardHand对象
        /// </summary>
        /// <param name="hai">手牌列表</param>
        /// <returns>如果牌列表是一个有效的组合，则返回一个pokerface.MsgCardHand对象，否则返回null</returns>
        public static pokerface.MsgCardHand  agariConvertMsgCardHand(int[] hai)
        {
            var key = calcKey(hai);
            if (!agariTable.ContainsKey(key))
            {
                return null;
            }

            var agari = agariTable[key];
            var ct = (pokerface.CardHandType)(agari & 0xff);

            var msgCardhand = new pokerface.MsgCardHand();
            msgCardhand.cardHandType = (int)ct;

            // 排序，让大的牌在前面
            Array.Sort(hai, (x, y) =>
            {
                return y - x;
            });

            // 如果是顺子类型则需要检查
            int flushLength = ((int)agari >> 16) & 0xff;
            if (flushLength > 0)
            {
                if (!agariFlushVerify(ct, flushLength))
                {
                    return null;
                }
            }

            var cardsNew = new List<int>();
            switch (ct)
            {
                case pokerface.CardHandType.TripletPair:
                case pokerface.CardHandType.Triplet2X2Pair:
                    // 确保3张在前面，对子在后面
                    for (var i = 0; i < hai.Length; i++)
                    {
                        var h = hai[i];
                        if (slots[h/4] == 3)
                        {
                            cardsNew.Add(h);
                        }
                    }
                    for (var i = 0; i < hai.Length; i++)
                    {
                        var h = hai[i];
                        if (slots[h / 4] != 3)
                        {
                            cardsNew.Add(h);
                        }
                    }
                    break;
                default:
                    cardsNew.AddRange(hai);
                    break;
            }

            msgCardhand.cards.AddRange(cardsNew);

            if (ct == pokerface.CardHandType.Triplet)
            {
                if (msgCardhand.cards[0] / 4 == (int)pokerface.CardID.R3H / 4)
                {
                    // 如果是3个3，而且不包含红桃3，则把牌组改为炸弹，而不是三张
                    var foundR3H = false;
                    foreach (var c in msgCardhand.cards)
                    {
                        if (c == (int)pokerface.CardID.R3H)
                        {
                            foundR3H = true;
                            break;
                        }
                    }

                    if (!foundR3H)
                    {
                        msgCardhand.cardHandType = (int)pokerface.CardHandType.Bomb;
                    }
                }
                else if (msgCardhand.cards[0] / 4 == (int)pokerface.CardID.AH / 4)
                {
                    // 3张A也是炸弹
                    msgCardhand.cardHandType = (int)pokerface.CardHandType.Bomb;
                }
            }
            return msgCardhand;
        }

        private static bool agariFlushVerify(CardHandType ct, int flushLength)
        {
            int flushElmentCount = 0;
            switch (ct)
            {
                case CardHandType.Flush:
                    flushElmentCount = 1;
                    break;
                case CardHandType.Pair2X:
                    flushElmentCount = 2;
                    break;
                case CardHandType.Triplet2X:
                case CardHandType.Triplet2X2Pair:
                    flushElmentCount = 3;
                    break;
            }

            if (flushElmentCount == 0)
            {
                // 不是顺子类型
                return true;
            }

            int flushBegin = 0;
            
            // 跳过2
            for(int i = 1; i < slots.Length; i++)
            {
                if (slots[i] == flushElmentCount)
                {
                    flushBegin = i;
                    break;
                }
            }

            int flushLengthVerify = 1;
            for(int i = flushBegin+1; i < slots.Length; i++)
            {
                if (slots[i] != flushElmentCount)
                {
                    break;
                }
                flushLengthVerify++;
            }

            return flushLengthVerify >= flushLength;
        }

        /// <summary>
        /// 判断当前的手牌是否大于上一手牌
        /// </summary>
        /// <param name="prevCardHand">上一手牌</param>
        /// <param name="current">当前的手牌</param>
        /// <returns>如果大于则返回true，其他各种情形都会返回false</returns>
        public static bool agariGreatThan(pokerface.MsgCardHand prevCardHand, pokerface.MsgCardHand current)
        {
            // 如果当前的是炸弹
            if (current.cardHandType == (int)pokerface.CardHandType.Bomb)
            {
                // 上一手不是炸弹
                if (prevCardHand.cardHandType != (int)pokerface.CardHandType.Bomb)
                {
                    return true;
                }

                // 上一手也是炸弹，则比较炸弹牌的大小，大丰关张不存在多于4个牌的炸弹
                return current.cards[0]/4 > prevCardHand.cards[0]/4;
            }

            // 如果上一手牌是炸弹
            if (prevCardHand.cardHandType == (int)pokerface.CardHandType.Bomb)
            {
                return false;
            }

            // 必须类型匹配
            if (prevCardHand.cardHandType != current.cardHandType)
            {
                return false;
            }

            // 张数匹配
            if (prevCardHand.cards.Count != current.cards.Count)
            {
                return false;
            }

            // 单张时，2是最大的
            if (prevCardHand.cardHandType == (int)pokerface.CardHandType.Single)
            {
                if (prevCardHand.cards[0] / 4 == 0)
                {
                    return false;
        
                }

                if (current.cards[0] / 4 == 0)
                {
                    return true;
                }
            }

            // 现在只比较最大牌的大小
            return current.cards[0]/4 > prevCardHand.cards[0]/4;
        }

        /// <summary>
        /// 寻找比上一手牌大的所有有效牌组
        /// 这个主要用于自动打牌以及给出提示之类
        /// </summary>
        /// <param name="prev">上一手牌</param>
        /// <param name="hands">当前手上所有的牌</param>
        /// <returns>返回一个牌组列表，如果没有有效牌组，该列表长度为0</returns>
        public static List<pokerface.MsgCardHand> FindGreatThanCardHand(pokerface.MsgCardHand prev, List<int> hands, int specialCardID)
        {
            var prevCT = (pokerface.CardHandType)prev.cardHandType;
            bool isBomb = false;
            List<pokerface.MsgCardHand> tt = null;
            if (specialCardID >= 0)
            {
                tt = new List<MsgCardHand>();
                var cardHand = new MsgCardHand();
                cardHand.cardHandType = (int)pokerface.CardHandType.Single;
                cardHand.cards.AddRange(extractCardByRank(hands, 0, 1));
                tt.Add(cardHand);
                return tt;
            }

            switch (prevCT)
            {
                case pokerface.CardHandType.Bomb:
                    tt = FindBombGreatThan(prev, hands);
                    isBomb = true;
                    break;
                case pokerface.CardHandType.Flush:
                    tt = FindFlushGreatThan(prev, hands);
                    break;
                case pokerface.CardHandType.Single:
                    tt = FindSingleGreatThan(prev, hands);
                    break;
                case pokerface.CardHandType.Pair:
                    tt = FindPairGreatThan(prev, hands);
                    break;
                case pokerface.CardHandType.Pair2X:
                    tt = FindPair2XGreatThan(prev, hands);
                    break;
                case pokerface.CardHandType.Triplet:
                    tt = FindTripletGreatThan(prev, hands);
                    break;
                case pokerface.CardHandType.Triplet2X:
                    tt = FindTriplet2XGreatThan(prev, hands);
                    break;
                case pokerface.CardHandType.Triplet2X2Pair:
                    tt = FindTriplet2X2PairGreatThan(prev, hands);
                    break;
                case pokerface.CardHandType.TripletPair:
                    tt = FindTripletPairGreatThan(prev, hands);
                    break;
            }

            if (!isBomb)
            {
                var tt2 = FindBomb(hands);
                tt.AddRange(tt2);
            }
            return tt;
        }

        /// <summary>
        /// 寻找手牌上的所有炸弹
        /// </summary>
        /// <param name="hands">手上的所有牌</param>
        /// <returns>返回一个牌组列表，如果没有有效牌组，该列表长度为0</returns>
        private static List<MsgCardHand> FindBomb(List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();
            ResetSlots(hands);

            for (var newBombSuitID = 0; newBombSuitID < (int)pokerface.CardID.AH / 4; newBombSuitID++)
            {
                if (slots[newBombSuitID] > 3)
                {
                    var cardHand = new MsgCardHand();
                    cardHand.cardHandType = (int)pokerface.CardHandType.Bomb;
                    cardHand.cards.AddRange(extractCardByRank(hands, newBombSuitID, 4));
                    cardHands.Add(cardHand);
                }
            }

            // 如果有3个ACE，也是炸弹
            if (slots[(int)pokerface.CardID.AH / 4] > 2)
            {
                var cardHand = new MsgCardHand();
                cardHand.cardHandType = (int)pokerface.CardHandType.Bomb;
                cardHand.cards.AddRange(extractCardByRank(hands, (int)pokerface.CardID.AH / 4, 3));
                cardHands.Add(cardHand);
            }


            // 黑桃梅花方块3组成炸弹
            List<int> three = new List<int>();
            foreach (var h in hands)
            {
                if (h / 4 == (int)pokerface.CardID.R3H / 4 && h != (int)pokerface.CardID.R3H)
                {
                    three.Add(h);
                }
            }

            if (three.Count == 3)
            {
                var cardHand = new MsgCardHand();
                cardHand.cardHandType = (int)pokerface.CardHandType.Bomb;
                cardHand.cards.AddRange(three);
                cardHands.Add(cardHand);
            }

            return cardHands;
        }
        /// <summary>
        /// 寻找所有大于上一手"连3张+两对子"的有效组合
        /// </summary>
        /// <param name="prev">上一手牌</param>
        /// <param name="hands">手上的所有牌</param>
        /// <returns>返回一个牌组列表，如果没有有效牌组，该列表长度为0</returns>
        private static List<MsgCardHand> FindTriplet2X2PairGreatThan(MsgCardHand prev, List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();
            ResetSlots(hands);

            var pairLength = 4;

            if (prev.cards.Count > 10) {
                pairLength = 6;
            }

            var flushLen = prev.cards.Count - pairLength;// 减去N个对子
            var bombCardRankID = prev.cards[0] / 4;
            var seqLength = flushLen / 3;
            for (var newBombSuitID = bombCardRankID + 1; newBombSuitID <= (int)pokerface.CardID.AH / 4;)
            {
                var testBombRankID = newBombSuitID;
                var found = true;
                for (var i = 0; i < seqLength; i++)
                {
                    if (slots[testBombRankID - i] < 3)
                    {
                        newBombSuitID = newBombSuitID + 1;

                        found = false;

                        break;

                    }
                }

                // 找到了
                if (found)
                {

                    var left = newBombSuitID + 1 - seqLength;
                    var right = newBombSuitID;

                    var pairCount = 0;
                    List<int> pairAble = new List<int>();
                    for (var testPair = 0; testPair < left; testPair++)
                    {
                        if (slots[testPair] > 1)
                        {
                            pairCount++;
                            pairAble.Add(testPair);
                        }
                    }

                    for (var testPair = right + 1; testPair < (int)pokerface.CardID.JOB / 4; testPair++)
                    {
                        if (slots[testPair] > 1)
                        {
                            pairCount++;
                            pairAble.Add(testPair);
                        }
                    }

                    if (pairCount >= seqLength)
                    {
                        // 此处不在遍历各种对子组合
                        var cardHand = new MsgCardHand();
                        cardHand.cardHandType = (int)pokerface.CardHandType.Triplet2X2Pair;
                        cardHand.cards.AddRange(extractCardByRanks(hands, left, right, 3));
                        var addPairCount = 0;
                        foreach (var pp in pairAble)
                        {
                            cardHand.cards.AddRange(extractCardByRank(hands, pp, 2));
                            addPairCount++;
                            if (addPairCount == seqLength)
                            {
                                break;
                            }
                        }

                        cardHands.Add(cardHand);
                    }

                    newBombSuitID = newBombSuitID + 1;
                }
            }

            return cardHands;
        }
        /// <summary>
        /// 寻找所有大于上一手"3张+对子"的有效组合
        /// </summary>
        /// <param name="prev">上一手牌</param>
        /// <param name="hands">手上的所有牌</param>
        /// <returns>返回一个牌组列表，如果没有有效牌组，该列表长度为0</returns>
        private static List<MsgCardHand> FindTripletPairGreatThan(MsgCardHand prev, List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();
            ResetSlots(hands);

            var flushLen = prev.cards.Count - 2;// 减去对子
            var bombCardRankID = prev.cards[0] / 4;
            var seqLength = flushLen / 3;
            for (var newBombSuitID = bombCardRankID + 1; newBombSuitID <= (int)pokerface.CardID.AH / 4;)
            {
                var testBombRankID = newBombSuitID;
                var found = true;
                for (var i = 0; i < seqLength; i++)
                {
                    if (slots[testBombRankID - i] < 3)
                    {
                        newBombSuitID = newBombSuitID + 1;

                        found = false;

                        break;

                    }
                }

                // 找到了
                if (found)
                {

                    var left = newBombSuitID + 1 - seqLength;
                    var right = newBombSuitID;

                    var pairCount = 0;
                    List<int> pairAble = new List<int>();
                    for (var testPair = 0; testPair < left; testPair++)
                    {
                        if (slots[testPair] > 1)
                        {
                            pairCount++;
                            pairAble.Add(testPair);
                        }
                    }

                    for (var testPair = right+1; testPair < (int)pokerface.CardID.JOB / 4; testPair++)
                    {
                        if (slots[testPair] > 1)
                        {
                            pairCount++;
                            pairAble.Add(testPair);
                        }
                    }

                    if (pairCount > 0)
                    {
                        // 此处不再遍历各个对子
                        var cardHand = new MsgCardHand();
                        cardHand.cardHandType = (int)pokerface.CardHandType.TripletPair;
                        cardHand.cards.AddRange(extractCardByRank(hands, left, 3));
                        cardHand.cards.AddRange(extractCardByRank(hands, pairAble[0], 2));
                        cardHands.Add(cardHand);
                    }

                    newBombSuitID = newBombSuitID + 1;

                }
            }

            return cardHands;
        }
        /// <summary>
        /// 寻找所有大于上一手"3张"的有效组合
        /// </summary>
        /// <param name="prev">上一手牌</param>
        /// <param name="hands">手上的所有牌</param>
        /// <returns>返回一个牌组列表，如果没有有效牌组，该列表长度为0</returns>
        private static List<MsgCardHand> FindTripletGreatThan(MsgCardHand prev, List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();
            ResetSlots(hands);

            var bombCardRankID = prev.cards[0] / 4;

            // 找一个较大的三张
            for (var newBombSuitID = bombCardRankID + 1; newBombSuitID < (int)pokerface.CardID.JOB / 4; newBombSuitID++)
            {
                if (slots[newBombSuitID] > 2)
                {
                    var cardHand = new MsgCardHand();
                    cardHand.cardHandType = (int)pokerface.CardHandType.Triplet;
                    cardHand.cards.AddRange(extractCardByRank(hands, newBombSuitID, 3));
                    cardHands.Add(cardHand);
                }
            }

            return cardHands;
        }
        /// <summary>
        /// 寻找所有大于上一手"连3张"的有效组合
        /// </summary>
        /// <param name="prev">上一手牌</param>
        /// <param name="hands">手上的所有牌</param>
        /// <returns>返回一个牌组列表，如果没有有效牌组，该列表长度为0</returns>
        private static List<MsgCardHand> FindTriplet2XGreatThan(MsgCardHand prev, List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();

            ResetSlots(hands);

            var flushLen = prev.cards.Count;
            var bombCardRankID = prev.cards[0] / 4; // 最大的顺子牌rank
            var seqLength = flushLen /3;
            for (var newBombSuitID = bombCardRankID + 1; newBombSuitID <= (int)pokerface.CardID.AH / 4;)
            {
                var testBombRankID = newBombSuitID;
                var found = true;
                for (var i = 0; i < seqLength; i++)
                {
                    if (slots[testBombRankID - i] < 3)
                    {
                        newBombSuitID = newBombSuitID + 1;

                        found = false;

                        break;

                    }
                }

                // 找到了
                if (found)
                {
                    var cardHand = new MsgCardHand();
                    cardHand.cardHandType = (int)pokerface.CardHandType.Triplet2X;
                    cardHand.cards.AddRange(extractCardByRanks(hands, testBombRankID - seqLength + 1, testBombRankID, 3));
                    cardHands.Add(cardHand);

                    newBombSuitID = newBombSuitID + 1;
                }

                
            }

            return cardHands;
        }
        /// <summary>
        /// 寻找所有大于上一手"连对"的有效组合
        /// </summary>
        /// <param name="prev">上一手牌</param>
        /// <param name="hands">手上的所有牌</param>
        /// <returns>返回一个牌组列表，如果没有有效牌组，该列表长度为0</returns>
        private static List<MsgCardHand> FindPair2XGreatThan(MsgCardHand prev, List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();

            ResetSlots(hands);

            var flushLen = prev.cards.Count;
            var bombCardRankID = prev.cards[0] / 4; // 最大的顺子牌rank
            var seqLength = flushLen / 2;
            for (var newBombSuitID = bombCardRankID + 1; newBombSuitID <= (int)pokerface.CardID.AH / 4;)
            {
                var testBombRankID = newBombSuitID;
                var found = true;
                for (var i = 0; i < seqLength; i++)
                {
                    if (slots[testBombRankID - i] < 2)
                    {
                        newBombSuitID = newBombSuitID + 1;

                        found = false;

                        break;

                    }
                }

                // 找到了
                if (found)
                {
                    var cardHand = new MsgCardHand();
                    cardHand.cardHandType = (int)pokerface.CardHandType.Pair2X;
                    cardHand.cards.AddRange(extractCardByRanks(hands, testBombRankID - seqLength + 1, testBombRankID, 2));
                    cardHands.Add(cardHand);

                    newBombSuitID = newBombSuitID + 1;
                }

               
            }

            return cardHands;
        }
        /// <summary>
        /// 寻找所有大于上一手"对子"的有效组合
        /// </summary>
        /// <param name="prev">上一手牌</param>
        /// <param name="hands">手上的所有牌</param>
        /// <returns>返回一个牌组列表，如果没有有效牌组，该列表长度为0</returns>
        private static List<MsgCardHand> FindPairGreatThan(MsgCardHand prev, List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();
            ResetSlots(hands);

            var bombCardRankID = prev.cards[0] / 4;

            // 找一个较大的对子
            for (var newBombSuitID = bombCardRankID + 1; newBombSuitID < (int)pokerface.CardID.JOB / 4; newBombSuitID++)
            {
                if (slots[newBombSuitID] > 1)
                {
                    var cardHand = new MsgCardHand();
                    cardHand.cardHandType = (int)pokerface.CardHandType.Pair;
                    cardHand.cards.AddRange(extractCardByRank(hands, newBombSuitID, 2));
                    cardHands.Add(cardHand);
                }
            }

            return cardHands;
        }
        /// <summary>
        /// 寻找所有大于上一手"单张"的有效组合
        /// </summary>
        /// <param name="prev">上一手牌</param>
        /// <param name="hands">手上的所有牌</param>
        /// <returns>返回一个牌组列表，如果没有有效牌组，该列表长度为0</returns>
        private static List<MsgCardHand> FindSingleGreatThan(MsgCardHand prev, List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();
            ResetSlots(hands);

            var bombCardRankID = prev.cards[0] / 4;
            if (bombCardRankID == 0)
            {
                // 2已经是最大的单张了
                return cardHands;
            }

            // 找一个较大的单张
            for (var newBombSuitID = bombCardRankID + 1; newBombSuitID < (int)pokerface.CardID.JOB / 4; newBombSuitID++)
            {
                if (slots[newBombSuitID] > 0)
                {
                    var cardHand = new MsgCardHand();
                    cardHand.cardHandType = (int)pokerface.CardHandType.Single;
                    cardHand.cards.AddRange(extractCardByRank(hands, newBombSuitID, 1));
                    cardHands.Add(cardHand);
                }
            }

            // 自己有2，那就是最大
            if (slots[0] > 0)
            {
                var cardHand = new MsgCardHand();
                cardHand.cardHandType = (int)pokerface.CardHandType.Single;
                cardHand.cards.AddRange(extractCardByRank(hands, 0, 1));
                cardHands.Add(cardHand);
            }

            return cardHands;
        }

        /// <summary>
        /// 寻找所有大于上一手"顺子"的有效组合
        /// </summary>
        /// <param name="prev">上一手牌</param>
        /// <param name="hands">手上的所有牌</param>
        /// <returns>返回一个牌组列表，如果没有有效牌组，该列表长度为0</returns>
        private static List<MsgCardHand> FindFlushGreatThan(MsgCardHand prev, List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();

            ResetSlots(hands);

            var flushLen = prev.cards.Count;
            var bombCardRankID = prev.cards[0] / 4; // 最大的顺子牌rank
            var seqLength = flushLen / 1;
            for (var newBombSuitID = bombCardRankID+1; newBombSuitID <= (int)pokerface.CardID.AH/4 ;)
            {
                var testBombRankID = newBombSuitID;
                var found = true;
                for (var i = 0; i < seqLength; i++)
                {
                    if (slots[testBombRankID - i] < 1)
                    {
                        newBombSuitID = newBombSuitID + 1;

                        found = false;

                        break;
        
                    }
                }

                // 找到了
                if (found)
                {
                    var cardHand = new MsgCardHand();
                    cardHand.cardHandType = (int)pokerface.CardHandType.Flush;
                    cardHand.cards.AddRange(extractCardByRanks(hands, testBombRankID- seqLength + 1, testBombRankID, 1));
                    cardHands.Add(cardHand);

                    newBombSuitID = newBombSuitID + 1;
                }

               
            }

            return cardHands;
        }
        /// <summary>
        /// 寻找所有大于上一手"炸弹"的有效组合
        /// </summary>
        /// <param name="prev">上一手牌</param>
        /// <param name="hands">手上的所有牌</param>
        /// <returns>返回一个牌组列表，如果没有有效牌组，该列表长度为0</returns>
        private static List<MsgCardHand> FindBombGreatThan(MsgCardHand prev, List<int> hands)
        {
            // 注意不需要考虑333这种炸弹，因为他是最小的，而现在是寻找一个大于某个炸弹的炸弹
            List<MsgCardHand> cardHands = new List<MsgCardHand>();

            ResetSlots(hands);

            var bombCardRankID = (prev.cards[0]) / 4;
            for (var newBombSuitID = bombCardRankID + 1; newBombSuitID < (int)pokerface.CardID.AH / 4; newBombSuitID++)
            {
                if (slots[newBombSuitID] > 3)
                {
                    var cardHand = new MsgCardHand();
                    cardHand.cardHandType = (int)pokerface.CardHandType.Bomb;
                    cardHand.cards.AddRange(extractCardByRank(hands, newBombSuitID, 4));
                    cardHands.Add(cardHand);
                }
            }

            // 如果有3个ACE，也是炸弹
            if (slots[(int)pokerface.CardID.AH/4] > 2)
            {
                var cardHand = new MsgCardHand();
                cardHand.cardHandType = (int)pokerface.CardHandType.Bomb;
                cardHand.cards.AddRange(extractCardByRank(hands, (int)pokerface.CardID.AH / 4, 3));

                cardHands.Add(cardHand);
            }

            return cardHands;
        }
        /// <summary>
        /// 根据手牌列表填充slots用于查找各种牌
        /// </summary>
        /// <param name="hands">手牌列表</param>
        private static void ResetSlots(List<int> hands)
        {
            for (var i = 0; i < slots.Length; i++)
            {
                slots[i] = 0;
            }

            foreach (var h in hands)
            {
                slots[h / 4]++;
            }
        }
        /// <summary>
        /// 根据rank，从手牌上提取若干张该rank的牌
        /// </summary>
        /// <param name="hands">手牌列表</param>
        /// <param name="rank">rank值</param>
        /// <param name="count">提取张数</param>
        /// <returns></returns>
        static List<int> extractCardByRank(List<int>hands, int rank, int count)
        {
            var extract = new List<int>();
            foreach(var h in hands)
            {
                if (h/4 == rank)
                {
                    extract.Add(h);

                    if (extract.Count == count)
                    {
                        break;
                    }
                }

            }

            return extract;
        }
        /// <summary>
        /// 根据一个rank范围，提取位于该范围的所有牌，每一种牌提取若干张
        /// </summary>
        /// <param name="hands">手牌列表</param>
        /// <param name="rankStart">起始rank</param>
        /// <param name="rankStop">最大的rank</param>
        /// <param name="countEach">每一个rank提取多少张</param>
        /// <returns></returns>
        static List<int> extractCardByRanks(List<int> hands, int rankStart, int rankStop, int countEach)
        {
            var extract = new List<int>();
            for (var rank = rankStart; rank <= rankStop; rank++)
            {
                var ce = 0;
                foreach (var h in hands)
                {
                    if (h / 4 == rank)
                    {
                        extract.Add(h);
                        ce++;
                        if (ce == countEach)
                        {
                            break;
                        }
                    }
                }
            }

            return extract;
        }

        #region 本区仅仅用于自动化工具，客户端不要采纳
        public static pokerface.MsgCardHand SearchLongestDiscardCardHand(List<int> hands, int specialCardID)
        {
            hands.Sort();
            List<pokerface.MsgCardHand> cardHands = new List<MsgCardHand>();
            cardHands.AddRange(SearchLongestFlush(hands));

            cardHands.AddRange(SearchLongestPairX(hands));

            cardHands.AddRange(SearchLongestTriplet2XOrTriplet2X2Pair(hands));

            cardHands.AddRange(SearchUseableTripletOrTripletPair(hands));

            cardHands.AddRange(SearchBomb(hands));

            cardHands.AddRange(SearchUseableSingle(hands));

            cardHands.Sort((x, y) =>
            {
                return y.cards.Count - x.cards.Count;
            });

            var needR3h = specialCardID >= 0;
            if (needR3h)
            {
                foreach(var ch in cardHands)
                {
                    for (var i = 0; i < ch.cards.Count; i++)
                    {
                        if (ch.cards[i] == (int)pokerface.CardID.R3H)
                        {
                            return ch;
                        }
                    }
                }
            }

            return cardHands[0];
        }

        private static List<MsgCardHand> SearchUseableSingle(List<int> hands2)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();
            ResetSlots(hands2);

            // 找一个较大的单张
            for (var newBombSuitID = 1; newBombSuitID < (int)pokerface.CardID.JOB / 4; newBombSuitID++)
            {
                if (slots[newBombSuitID] > 0)
                {
                    var cardHand = new MsgCardHand();
                    cardHand.cardHandType = (int)pokerface.CardHandType.Single;
                    cardHand.cards.AddRange(extractCardByRank(hands2, newBombSuitID, 1));
                    cardHands.Add(cardHand);
                }
            }

            // 自己有2，那就是最大
            if (slots[0] > 0)
            {
                var cardHand = new MsgCardHand();
                cardHand.cardHandType = (int)pokerface.CardHandType.Single;
                cardHand.cards.AddRange(extractCardByRank(hands2, 0, 1));
                cardHands.Add(cardHand);
            }

            return cardHands;
        }

        private static List<MsgCardHand> SearchBomb(List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();
            ResetSlots(hands);

            for (var newBombSuitID = 0; newBombSuitID < (int)pokerface.CardID.AH / 4; newBombSuitID++)
            {
                if (slots[newBombSuitID] > 3)
                {
                    var cardHand = new MsgCardHand();
                    cardHand.cardHandType = (int)pokerface.CardHandType.Bomb;
                    cardHand.cards.AddRange(extractCardByRank(hands, newBombSuitID, 4));
                    cardHands.Add(cardHand);
                }
            }

            // 如果有3个ACE，也是炸弹
            if (slots[(int)pokerface.CardID.AH / 4] > 2)
            {
                var cardHand = new MsgCardHand();
                cardHand.cardHandType = (int)pokerface.CardHandType.Bomb;
                cardHand.cards.AddRange(extractCardByRank(hands, (int)pokerface.CardID.AH / 4, 3));
                cardHands.Add(cardHand);
            }

            // 黑桃梅花方块3组成炸弹
            foreach (var h in hands)
            {
                List<int> three = new List<int>();
                if (h / 4 == (int)pokerface.CardID.R3H / 4 && h != (int)pokerface.CardID.R3H)
                {
                    three.Add(h);
                }

                if (three.Count == 3)
                {
                    var cardHand = new MsgCardHand();
                    cardHand.cardHandType = (int)pokerface.CardHandType.Bomb;
                    cardHand.cards.AddRange(three);
                    cardHands.Add(cardHand);
                }
            }

            return cardHands;
        }

        private static List<MsgCardHand> SearchUseableTripletOrTripletPair(List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();
            ResetSlots(hands);

            for (var newBombSuitID = 0; newBombSuitID <= (int)pokerface.CardID.AH / 4;)
            {
                var testBombRankID = newBombSuitID;
                var found = true;
                for (var i = 0; i < 1; i++)
                {
                    if (slots[testBombRankID - i] < 3)
                    {
                        newBombSuitID = newBombSuitID + 1;

                        found = false;

                        break;

                    }
                }

                // 找到了
                if (found)
                {
                    var cardHand = new MsgCardHand();
                    cardHand.cards.AddRange(extractCardByRank(hands, newBombSuitID, 3));
                    cardHands.Add(cardHand);

                    var left = newBombSuitID ;
                    var right = newBombSuitID;

                    var pairCount = 0;
                    List<int> pairAble = new List<int>();
                    for (var testPair = 0; testPair < left; testPair++)
                    {
                        if (slots[testPair] > 1)
                        {
                            pairCount++;
                            pairAble.Add(testPair);
                        }
                    }

                    for (var testPair = right + 1; testPair < (int)pokerface.CardID.JOB / 4; testPair++)
                    {
                        if (slots[testPair] > 1)
                        {
                            pairCount++;
                            pairAble.Add(testPair);
                        }
                    }

                    if (pairCount > 0)
                    {
                        // 此处不再遍历各个对子
                        cardHand.cards.AddRange(extractCardByRank(hands, pairAble[0], 2));
                        
                    }

                    newBombSuitID = newBombSuitID + 1;

                }
            }

            return cardHands;
        }

        private static List<MsgCardHand> SearchLongestTriplet2XOrTriplet2X2Pair(List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();

            ResetSlots(hands);

            for (var newBombSuitID = 0; newBombSuitID <= (int)pokerface.CardID.AH / 4;)
            {
                var testBombRankID = newBombSuitID;
                for (var i = 0; i < 13; i++)
                {
                    if (slots[testBombRankID + i] < 3)
                    {

                        // 找到了
                        if (i >= 2)
                        {
                            var cardHand = new MsgCardHand();
                            cardHand.cards.AddRange(extractCardByRanks(hands, testBombRankID, testBombRankID + i - 1, 3));
                            cardHands.Add(cardHand);

                            // 寻找2个对子
                            var left = testBombRankID;
                            var right = testBombRankID + i - 1;

                            var pairCount = 0;
                            List<int> pairAble = new List<int>();
                            for (var testPair = 0; testPair < left; testPair++)
                            {
                                if (slots[testPair] > 1)
                                {
                                    pairCount++;
                                    pairAble.Add(testPair);
                                }
                            }

                            for (var testPair = right + 1; testPair < (int)pokerface.CardID.JOB / 4; testPair++)
                            {
                                if (slots[testPair] > 1)
                                {
                                    pairCount++;
                                    pairAble.Add(testPair);
                                }
                            }

                            if (pairCount >= i)
                            {
                                // 此处不在遍历各种对子组合
                                var addPairCount = 0;
                                foreach (var pp in pairAble)
                                {
                                    cardHand.cards.AddRange(extractCardByRank(hands, pp, 2));
                                    addPairCount++;
                                    if (addPairCount == i)
                                    {
                                        break;
                                    }
                                }
                                
                            }
                        }


                        break;

                    }
                }



                newBombSuitID = newBombSuitID + 1;
            }

            return cardHands;
        }

        private static List<MsgCardHand> SearchLongestPairX(List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();

            ResetSlots(hands);

            for (var newBombSuitID = 0; newBombSuitID <= (int)pokerface.CardID.AH / 4;)
            {
                var testBombRankID = newBombSuitID;
                for (var i = 0; i < 13; i++)
                {
                    if (slots[testBombRankID + i] < 2)
                    {

                        // 找到了
                        if (i >= 1)
                        {
                            var cardHand = new MsgCardHand();
                            cardHand.cards.AddRange(extractCardByRanks(hands, testBombRankID, testBombRankID + i - 1, 2));
                            cardHands.Add(cardHand);
                        }

                        break;

                    }
                }

                newBombSuitID = newBombSuitID + 1;
            }


            return cardHands;
        }

        private static List<MsgCardHand> SearchLongestFlush(List<int> hands)
        {
            List<MsgCardHand> cardHands = new List<MsgCardHand>();

            ResetSlots(hands);
            // 简单起见从3开始搜索，不考虑ACE开始的类似12345这种
            for (var newBombSuitID = 1; newBombSuitID <= (int)pokerface.CardID.AH / 4;)
            {
                var testBombRankID = newBombSuitID;
                for (var i = 0; i < 13; i++)
                {
                    if (slots[testBombRankID + i] < 1)
                    {

                        // 找到了
                        if (i > 4)
                        {
                            var cardHand = new MsgCardHand();
                            cardHand.cards.AddRange(extractCardByRanks(hands, testBombRankID , testBombRankID+i-1, 1));
                            cardHands.Add(cardHand);
                        }

                        break;
                    }
                }

                newBombSuitID = newBombSuitID + 1;
            }

            return cardHands;
        }

        #endregion
        static AgariIndex()
        {
            agariTable[0x423a35c7] = 0xa0a01;
            agariTable[0x14d] = 0x30908;
            agariTable[0x8235] = 0x50f08;
            agariTable[0x4] = 0x402;
            agariTable[0x1] = 0x103;
            agariTable[0x2] = 0x204;
            agariTable[0x20] = 0x507;
            agariTable[0x10f447] = 0x70701;
            agariTable[0xde] = 0x20405;
            agariTable[0xcfa] = 0x20a09;
            agariTable[0x21] = 0x20608;
            agariTable[0xd05] = 0x40c08;
            agariTable[0xa98ac7] = 0x80801;
            agariTable[0x69f6bc7] = 0x90901;
            agariTable[0x16] = 0x20605;
            agariTable[0x8ae] = 0x30605;
            agariTable[0x56ce] = 0x40805;
            agariTable[0x3640e] = 0x50a05;
            agariTable[0x21e88e] = 0x60c05;
            agariTable[0x153158e] = 0x70e05;
            agariTable[0x3] = 0x306;
            agariTable[0x2b67] = 0x50501;
            agariTable[0x1b207] = 0x60601;
            agariTable[0x2964619c7] = 0xb0b01;
            agariTable[0x19debd01c7] = 0xc0c01;
            agariTable[0x515a6] = 0x30f09;
        }
    }
}
