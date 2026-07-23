| # | Name | Pattern | Date | Time | Solved unaided? | Notes |
**Week 1 — Arrays & Hashing**
1   |Two Sum| use hash map | 21.04 | Start time: 12:26, End time: 12:33 | Unaided | can use Map() here (hash.set, has, get)
2   |Contains Duplicate| use hash set | 21.04 | Start time: 12:39, End time: 12:44 | Unaided | can use Set() here (set.add, has, delete)
3   |Valid Anagram| used 2 maps and then compared to each other | 21.04 | Start time: 13:19, End time: 13:25 | Unaided | can use Map() here (hash.set, has, get)
4   |Group Anagram| sorted each string and put into map (sortedStr, array of matching strings) | 22.04 | Start time: 11:02, End time: 11:12 | Unaided | use split(""), join(""), sort(), push(), Array.from(map.values())
5   |Top K Frequent Elements| put num+freq into the map, sort by freq, return first k elements | 22.04 | Start time: 11:16, End time: 11:39 | Unaided for algorythm, needed help with JS syntax |     var sortedMap = new Map([...map.entries()].sort((a, b) => b[1] - a[1]));
    const firstK = new Map([...sortedMap].slice(0, k));

    return Array.from(firstK.keys());
6   |Product of Array Except Self | create prefix sums array which stores product of previous elements, 
and suffix sum array which stores product of next element, then for each element calculate product of prefixSums[i -1] * suffixSum[i + 1] | 23.04 | Start time: 11:52, End time: 12:10 | Unaided | Dont forget to reverse() suffixSums array
7.  |Valid Sudoku | Create 3 arrays with sets, for rows, cols and boxes. | 23.04 | Start time: 12:15, End time: 12:38 | Unaided for 70% but I needed help with var idx = Math.floor(row/3) * 3 + Math.floor(col/3); | DOnt forget if (curVal == '.') continue; and var idx = Math.floor(row/3) * 3 + Math.floor(col/3);
9.  | Longest Consecutive Sequence | use hash map(value+true). Then for each sequence start calculating its length and choose longest one(two pointers technique) | 24.04 | Start time: 12:05, End time: 12:30 | Aided for pattern+coding  |  if (!map.has(nums[i] - 1))  - means new sequence started. set length as 0 and while loop while (map.has(nums[i] + length)) to calculate curent sequence length. Then compare to previous lebngthes to find max one. RE-DO AGAIN LATER.

**Week 2 — Two Pointers & Sliding Window**
10  |Valid Palindrome| first clean string to be alphanumeric, then use two pointers at the start and the end and check if equal or not | 27.04 | Start time: 11:45, End time: 12:08 | Unaided for algorythm, had to check up alphanumeric check in JS | var isAlphanumeric = function(c) {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9');
}
11  |Two Sum II| It's binary search | 28.04 | Start time: 12:15 , End time: 12:23 | Aided with algorythm, solved unoptimal first | Use 2 pointers, at the start and the end, check its sum, if more or less, start++ or end--
12  |3Sum| Sort numbers, then use 2 pointers technique. Cur value + nums[left] + nums[right], move pointers depended on sum (sum < 0, sum > 0, sum == 0). Remove dups by moving pointers too by checking (nums[i-1], nums[left+1], nums[right-1]) | 29.04 | Start time: 11:15 , End time: 11:36 | Aided fully, forgot algorythm | for (var i = 0; i < nums.length - 2; i++) { //i is first element of triplet, so dont touch last 2 element RE-DO AGAIN LATER
if (i > 0 && nums[i] == nums[i-1]) continue; //remove duplicates
//again remove dups(similar for right)
                while (left < right && nums[left] == nums[left+1]) {
                    left++;
                }
13  |Container With Most Water| use only 2 pointers, move left if height[left] < height[right] and vice versa | 29.04 | Start time: 11:55 , End time: 12:15 | Aided with algorythm partially(when to move pointers) | RE-DO AGAIN LATER 
14  |Trapping Rain Water| track leftMax and rightMax and check if leftMax < rightMax or not | 29.04 | Start time: 18:35 , End time: 19:10 | Aided with algorythm(did about 30% by myself) | RE-DO AGAIN LATER 
15  |Best Time to Buy/Sell Stock| left = 0, right = 1. iterate through right one, one pass. SInce its alway one-directional and we cant return back, if we see right that is less than left, new left is that right (left = right) | 30.04 | Start time: 12:00 , End time: 12:15 | Aided with algorythm(did about 30% by myself) | RE-DO AGAIN LATER if you forgot
16  |Longest Substring Without Repeating Characters| two pointers and hash map | 30.04 | Start time: 12:00 , End time: 12:30 | Uniaded for 90% | 
17  |Longest Repeating Character Replacement|  | 01.05 | Start time: 15:35 , End time: 16:10 | Aided | RE-DO AGAIN LATER 
18  |Minimum Window Substring|  | 04.05 | Start time: 14:00 , End time: 14:58 | Aided | RE-DO AGAIN LATER 

**Week 3 — Stack & Binary Search**
19  |Valid Parentheses| Use stack + hash map | 04.05 | Start time: 16:35 , End time: 14:58 | Unaided for most part| -
20  |Min Stack| Store minimum elemnt alongside in stack itself - this.stack.push([val, minVal]); | 06.05 | Start time: 11:50 , End time: 12:02 | Aided | To init - this.stack = []; RE-DO
21  |Evaluate Reverse Polish Notation| Use stack to store result | 06.05 | Start time: 12:21 , End time: 12:39 | Unaided only with Math.trunc| res = Math.trunc(second / first); (dont use Math.floor it make -3.5 -4 and not -3 as expected)
!isNaN(cur) to check if integer
22  |Generate Parentheses| Use backtrack technique - recursion | 07.05 | Start time: 13:02 , End time: 14:02 | Aided | RE-DO
23  |Daily Temperatures| Monotonic stack pattern | 07.05 | Start time: 14:10 , End time: 14:40 | Aided | RE-DO
Fill n length array with 0: var res = new Array(n).fill(0);
24  |Car Fleet| Monotonic stack pattern | 08.05 | Start time: 22:25 , End time: 22:52 | Aided | RE-DO
var map = position.map((p, i) => [p, speed[i]]); - from array to map;
for (var [pos, speed] of map) - iterate map
25  |Binary Search| left, right, mid | 10.05 | Start time: 22:00 , End time: 22:05 | Unaided | while (left <= right) !
26  |Search 2D Matrix| Binary search, use while (left <= right), not iterate through rows & cols, to calculate current element: var curElement = matrix[Math.floor(mid / cols)][mid % cols]; | 10.05 | Start time: 22:12 , End time:  | Aided with algo | RE-DO; mid / cols = which row it is; mid % cols = which col it is;
27  |Koko Eating Bananas| Use binary search, left is 1h, right is max value in piles, find mid, for each mid calcultae time that would be spent to be eaten with such speed. Return left(its minimum time) | 11.05 | Start time: 12:45 , End time: 13:00 | Aided | RE-DO 
    function time(rate) {
        var sum = 0;
        for (var i = 0; i < piles.length; i++) {
            sum += Math.ceil(piles[i] / rate);
        }
        return sum;
    }
    //Maximum in array - var right = Math.max(...piles);
28  |Find Minimum in Rotated Sorted Array| Make early return if (nums[left] < nums[right]) return nums[left]; Search in the dip part. If one part is asc, skip it | 11.05 | Start time: 13:23 , End time: 13:40 | Aided | RE-DO (in search for the dip)
29  |Search in Rotated Sorted Array| Search in intervals between left and right ad chneh pointers accordingly | 11.05 | Start time: 14:48, End time:  | Unaided | 
left <= right — searching for an exact value (if testing is equal)
left < right — searching for a position/boundary
30  |Time Based Key-Value Store| use map, store [timestamp, value], res = curVals[mid][1]; left = mid + 1 | 12.05 | Start time: 15:35 , End time: 16:27 | Aided | Re-do
31  |Reverse Linked List| see code below, use newList = null, store rest of the list in tmp(will become curNode at athe end), add to curNode reversed list at this point, revList = curNode | 13.05 | Start time: 12:00 , End time: 12:17 | Aided partially | 
    while (curNode != null) {
        var tmp = curNode.next;
        curNode.next = newList;
        newList = curNode;
        curNode = tmp;
    }
32  |Merge Two Sorted Lists| Create newList and return newList.next; newNode = newList and iterates through cycle. Iterate through list1 and list2 and compare its values. Go next if either is appended. append list1 or list2 to newNode. Increment newNode = newNode.next | 13.05 | Start time: 12:20 , End time: 12:38 | Aided partially | Re-do, add remaining newNode.next = (list1 != null) ? list1 : list2;
33  |Reorder List| Two pointers (fast and slow) - use fast to reach middle of linked list. Reverse second half, then merge two lists | 13.05 | Start time: 20:20 , End time: 21:26 | Aided | Re-do, use while (second.next) when iterating last loop. var first = head; //this is pointers not lists
34  |Remove Nth Node From End of List| Use fast and slow pointers. First move fast pointer in for loop, i < n; Then move both slow and fast and connect| 16.05 | Start time: 15:35 , End time: 15:55 | Aided | Re-check
36  |Add Two Numbers| Use dummy node and carry = 0 (first digit in two digit number, carry = Math.floor(sum / 10);), curr = dummy. curr.next = new ListNode(sum % 10)| 16.05 | Start time: 22:48 , End time:  | Aided | Re-check 
37  |Linked List Cycle| Use fast and slow pinters, if fast pointer is EQUAL slow pointer, there is a cycle| 17.05 | Start time: 15:49 , End time: 15:58 | Unaided for 80% | Re-check 
38  |Find the Duplicate Number| Use binary search. We calculate mid index and then count all numbers, that are less or eqaul to that mid number. If there are more such numbers, it means duplicate is in that part, and vice versa. Low and high in this case are a set of potential numbers. | 17.05 | Start time: 17:13 , End time: 17:40 | Aided | Re-do
39  |LRU Cache| When I "get" the key, I first delete it, then set again, then return value. Thus, if its used, its always appended last to the list. So when it's time to find the LRU element, its always on the top of the list | 17.05 | Start time: 17:40 , End time: 17:56 | Aided | Re-check 
      this.cache.delete(this.cache.keys().next().value);  // keys().next().value returns first item's key
      //Init map and capacity:     this.cache = new Map(); this.capacity = capacity;
40  |Invert Binary Tree| Use recursion. invert left and right and then re-assign and retirn root | 18.05 | Start time: 10:49 , End time: 10:56 | Unaided for 50% | Re-check
41  |Maximum Depth of Binary Tree| Use recursion, return max of (maxDepth(left) || maxDepth(right))+1;| 18.05 | Start time: 10:57 , End time: 11:05 | Unaided for 50% | Re-check, if root == null, return 0;
42  |Diameter of Binary Tree| use DFS helper function and max = 0; | 19.05 | Start time: 09:04 , End time: 09:20 | Unaided for 70%| Re-check
43  |Balanced Binary Tree| use DFS - left, right, then check if either of them is -1; if Math.abs(left - right) > 1 return -1; height itself is Math.max(left, right) + 1 | 19.05 | Start time: 09:50 , End time: 10:16 | Aided | Re-check
Calculate height of a tree: return Math.max(left, right) + 1;
44  |Same Tree| Use recursion, same function. check if both null, then true, if either null then false, then check left and right and if vals are equal | 21.05 | Start time: 22:23 , End time: 22:26 | Unaided |
45  |Subtree of Another Tree| Use previous isSame tree check, to check if root and subroot are equal, and do recursion on isSubtree(root.left, subtree) || isSubtree(root.right, subtree) | 21.05 | Start time: 22:27 , End time: 22:41 | Unaided for 50% | Re-check
46  |Lowest Common Ancestor of a Binary Search Tree| BST (use while node != null and then move nodes to the right or left, comparing theit values). If both nodes are more than curVal- move right. If both are less = move left. Else (this means one is less and one is mor, or some is equal) curNode is Lowest Common Ancestor | 22.05 | Start time: 12:10 , End time: 12:22 | Aided | Re-do; var mode = root;
var lowestCommonAncestor = function(root, p, q) {
    var node = root;
    while (node != null) {
        if (p.val > node.val && q.val > node.val) {
            node = node.right;
        } else if (p.val < node.val && q.val < node.val) {
            node = node.left;
        } else {
            return node;
        }
    }
    return null;
};
47  |Binary Tree Level Order Traversal| Use BFS. start new curLevel = [] on each iteration. then iterate through queue and push to curLevel | 22.05 | Start time: 12:24 , End time: 12:31 | Aided for 70% | Re-check
function bfs(root) {
  if (!root) return;
  const queue = [root];
  while (queue.length) {
    const node = queue.shift();
    if (node.left) queue.push(node.left);
    if (node.right) queue.push(node.right);
  }
}
48  |Binary Tree Right Side View| Do level order traversal and push to the res last element of current level queue | 22.05 | Start time: 17:08 , End time: 17:12 | Unaided |if (i == curLen - 1) { res.push(curNode.val);}
49  |Count Good Nodes in Binary Tree| Use dfs. in dfs function send curMaxVal. curMaxVal = curNode.val and cnt++ if curNode.val >= curMaxVal. Return left+right+cnt| 25.05 | Start time: 19:38 , End time: 19:56 | Unaided | Re-check
50  |Validate Binary Search Tree| Use dfs(root, min(-Infinity), max(Infinity)). Check if root.val <= min || >= max, then false. Otherwise check dfs(root.left, min, root.val) && dfs(root.right, root.val, max) | 25.05 | Start time: 19:58 , End time: 20:21 | Aided for 80%| Re-check
51  |Kth Smallest Element in a BST| In-order traversal saves nodes in sorted order | 26.05 | Start time: 12:21 , End time:  | Aided | Re-do
var inorder = function(root, arr) {
    if (root == null) return arr;
    inorder(root.left, arr); - push smaller
    arr.push(root.val); - push current
    inorder(root.right, arr); - push larger
    return arr;
} if you want desc order, do push larger first then push smaller
52  |Construct Binary Tree from Preorder and Inorder Traversal| take preorder[0] as root, find it in inorder, recurse on the two halves. Use map to store val->i of inorder to lookup during recurse| 26.05 | Start time: 17:00 , End time: 17:29 |Aided| Re-do
Preorder is root → left → right, so preorder[0] is always the root of the current (sub)tree.
Inorder is left → root → right, so once you locate the root's position in inorder, everything to its left belongs to the left subtree, and everything to its right belongs to the right subtree.
var rootVal = preorder[preidx++]; //rootValue comes from preorder
        var newNode = new TreeNode(rootVal);
        var mid = indexMap.get(rootVal);
53  |Binary Tree Maximum Path Sum| set global maxGain. in maxGainPath(node) calculate - leftMaxGain, rightMaxGain(take max from comparing with 0). then calculate maxGain = max(maxGain, left+right+curVal). Return curVal + max(left, right) | 27.05 | Start time: 11:50 , End time: 12:20 | Aided | Re-check
dont name vars and functions with same name - will get error
54  |Serialize and Deserialize Binary Tree| Use dfs. concatenate with ",", then split by "," and send to str/new Node | 28.05 | Start time: 20:55 , End time: 21:34 | Aided | Re-check
str.split(",") for coverting string into array;
str.shift(); // deletes first element from array and returns it
55  |Implement Trie (Prefix Tree)| Store this.children = {} and bool isEnd(this will indicate if a full word or not). When inserting, foreach character create new Trie and proceed to the next character, set isEnd = true; | 29.05 | Start time: 17:31 , End time: 18:31 | Aided | Re-check
 let node = this; - to traverse
 if (!node.children[ch]) node.children[ch] = new Trie(); - if current word character is not in current node children
56  |Design Add and Search Words Data Structure| Use same technique for add and initiate as previous, but search is different | 29.05 | Start time: 18:33 , End time: 19:24 | Aided | Re-do
WordDictionary.prototype.search = function(word) {
    return this._search(word, 0);
};

WordDictionary.prototype._search = function(word, i) {
    if (word.length == i) return this.isEnd; //check is full word
    var ch = word[i];
    if (ch == '.') {
        for (var key in this.children) {
            if (this.children[key]._search(word, i + 1)) return true;
        }
        return false;
    }
    if (!this.children[ch]) return false;
    return this.children[ch]._search(word, i + 1);
};
57 |Word Search II| Build a trie from words, store whole word at then end nodes, then run dfs for each value in board| 01.06 | Start time: 12:00 , End time: 13:00 | Aided | Re-do
Build a trie:
    for (var word of words) {
        var node = root;
        for (var ch of word) {
            if (!node[ch]) node[ch] = {};
            node = node[ch];
        }
        node.word = word;
    }
58 |Kth Largest Element in a Stream| Use MinHeap | 02.06 | Start time: 11:55 , End time: 12:34 | Aided | Re-check - new syntax for heaps!!
this.heap = new MinPriorityQueue({ compare: (a, b) => a - b }); or new MinPriorityQueue()
this.heap.enqueue(num);
this.heap.size()
this.heap.dequeue();
this.heap.front();
59 |Last Stone Weight| Use MaxHeap, then check first two elements and enqueue again if needed, loop through while (heap.size() > 1) | 02.06 | Start time: 12:40 , End time: 13:00 | Unaided | can use var first = heap.dequeue(); directly
60 |K Closest Points to Origin| Use maxHeap | 04.06 | Start time: 11:44 , End time: 12:01 | Aided for syntax | Use max heap, because we dequeue all the larger ones, so smallest ones are left - const heap = new MaxPriorityQueue((p) => p[0]*p[0] + p[1]*p[1]);
61 |Kth Largest Element in an Array| Use MinPriorityQueue | 04.06 | Start time: 12:10, End time: 12:14 | Unaided | store in minheap, dequeue if size is more than k, the largest will be at the top 
62 |Task Scheduler| always run the task with the highest remaining count, and park tasks that just ran into a cooldown queue until they're available again | 05.06 | Start time: 12:30 , End time: 13:00 | Aided | RE-DO AGAIN 
//Create frequencies hash map:
    for (var task of tasks) {
        frequencies[task] = (frequencies[task] || 0) + 1;
    }
//Get value of object values
for (var freq of Object.values(frequencies))
63 |Design Twitter| store in this.users[userId] = { tweets: [], following: new Set() }; when get news feed, pull self tweets into max heap, then for each user that initial user is following, get their tweets too and put into same heap. Then retrieve first 10 tweets interatively. Tweets store as [tweetId, timestamp - global(timestamp++)]| 06.06 | Start time: 12:00, End time: 12:45 | Aided | Re-check
64 |Find Median from Data Stream| Use min and max heaps, Core idea: split the stream into two halves so the median is always at the boundary.
maxHeap holds the smaller half — its top is the largest of the small numbers.(e.g. 4,3,2)
minHeap holds the larger half — its top is the smallest of the large numbers.(e.g. 5,6,7)| 06.06 | Start time: 12:50, End time: 13:40 | Aided | Re-check
65 |Insert Interval| iterate through intervals and first add to res those intervals, who are less than new Interval. Then merge if overlapping. Then add remaining| 08.06 | Start time: 12:15 , End time: 12:30 | Aided for algorithm | Re-check
To merge, MODIFY newInterval in place, push newInterval after looping. 
        newInterval[0] = Math.min(intervals[i][0], newInterval[0]);
        newInterval[1] = Math.max(intervals[i][1], newInterval[1]);
66 |Merge Intervals| compare against and modify the last element of res | 08.06 | Start time: 21:35 , End time: 21:51 | Aided for 20% for algorithm | Re-check
67 |Non-overlapping Intervals| sort by second element ASC, keep prevEnd = intervals[0][1] and compare with current interval first element to find overlap | 10.06 | Start time: 12:45 , End time: 13:15 | Aided for 20% for algorithm | Re-check
68 |Non-overlapping Intervals| sort by second element ASC, keep prevEnd = intervals[0][1] and compare with current interval first element to find overlap | 10.06 | Start time: 12:45 , End time: 13:15 | Aided for 20% for algorithm | Re-check
69 |Minimum Interval to Include Each Query| for each query check min heap, which contains r - l + 1 of intervals that starts before query and right point. remove interva;s that ended before query. Put first element of the heap into res | 10.06 | Start time: 18:45 , End time:  | Aided | Re-do
FOR INTERVALS PROBLEMS ALWAYS SORT THEM FIRST
var c = queries.map((value, i) => [value, i]) //create map from array
var res = new Array(queries.length).fill(-1);
const heap = new MinPriorityQueue((el) => el[0]); // element is [size, end] → order by size
**Week 8 — Backtracking**
69 |Subsets| User recusrsion backtrack(start). push to res on every node, then run a loop: push to current path, backtrack(i+1), pop from path | 11.06 | Start time: 12:30 , End time: 13:00 | Aided | Re-check
70 |Subsets II| Sort and skip if duplicate, the rest is same as Subsets | 11.06 | Start time: 13:30 , End time: 14:00 | Unaided | if (i > start && nums[i] == nums[i - 1]) { continue }
71 |Combination Sum| Almost same as subsets. Make early returns when found target or curSum id > target | 12.06 | Start time: 16:30 , End time: 16:50 | Unaided for 70% | don't forget to push a copy of combination, not reference - res.push([...combination]); backtrack(i, curSum + candidates[i]); //i because we can re-use same number
i + 1 → each element used at most once (Subsets, Combinations, Combination Sum II).
i → each element reusable unlimited times (Combination Sum).
still passing "start" forward in both → prevents permutation duplicates like [2,3] and [3,2].
71 |Combination Sum II| Same as previous, but backtrack(i + 1, curSum), sort candidates at first, and check for duplicates | 12.06 | Start time: 16:55 , End time: 17:15 | Unaided for 80%| if (i > start && candidates[i] == candidates[i - 1]) continue; // skip duplicates
72 |Permutations| Same backtracking pattern. To push into res, check if path.length == nums.length. To iterate through all possible solutions, always start from var i = 0. To prevent duplicates, check if element already exists| 29.06 | Start time: 11:31 , End time: 11:53 | Aided for 30% | Re-check
        for (var i = 0; i < nums.length; i++) {
            if (path.includes(nums[i])) continue;
            path.push(nums[i]);
            backtrack(i + 1);
            path.pop();
        }
73 |Word Search| Use DFS instead. Backtracking here is to set board[r][c] = '#' to mark as visited and return it back | 30.06 | Start time: 12:03 , End time: 12:28 | Aided | Re-DO
you need to check if board[r][c] == word[curIndex]. If curIndex == word.length return true;
74 |Palindrome Partitioning| base case check if start == s.length; push to path only if palindrome. push to path substring from start to end. | 02.07 | Start time: 11:56 , End time: 12:24 | Aided | Re-DO
    function isPalindrome(l, r) {
        while (l < r) {
            if (s[l] != s[r]) return false;
            l++;
            r--;
        }
        return true;
    }
75 |Letter Combinations of a Phone Number| Make map, backtrack from 0; if idx == digits.length, push to res. then iterate every letter in mapped, using push +backtrack+ pop | 04.07 | Start time: 16:55 , End time:  17:11| Aided for 70%| Re-DO
    var map = {
        1 : "",
        2 : "abc",
        3 : "def",
        4 : "ghi",
        5 : "jkl",
        6 : "mno",
        7 : "pqrs",
        8 : "tuv",
        9 : "wxyz"
    };
76 |N-Queens| Keep track of cols = set(); posDiag = set, which contains row + col; negDiag = set, which contains row - col; Track queens positions in columns and iterate over columns| 05.07 | Start time: 12:27 , End time: 14:27 | Aided | Re-check
    function backtrack(row) {
        if (row == n) {
            res.push(queens.map(c => '.'.repeat(c) + 'Q' + '.'.repeat(n - c - 1)));
            return;
        }
        for (var col = 0; col < n; col++) {
            if (cols.has(col) || posDiag.has(row + col) || negDiag.has(row - col)) continue;
            queens[row] = col;
            cols.add(col); posDiag.add(row + col); negDiag.add(row - col);
            backtrack(row + 1);
            cols.delete(col); posDiag.delete(row + col); negDiag.delete(row - col);
        }
    }
77 |Number of Islands| Use DFS standard. Iterate over rows & cols and do dfs inside | 07.07 | Start time: 22:46 , End time: 00:20 | Aided | Re-do
78 |Clone Graph| Use DFS. before dfs create adjList = new Map(); inside dfs check if adjList has node.val and return if found. If not, create new clone Node(only val), set to node.val this cloneNode. Then iterate through node.neighbors and push to cloneNode.neighbors(dfs(neighbor)). Return cloneNode | 09.07 | Start time: 17:54 , End time: 18:08| Aided | Re-do
79 |Max Area of Island| for each rows and cols run dfs and maxArea = Math.max(maxArea, dfs(row, col)); Inside dfs, mark grid[r][c] == 0, so it wont return there. var curSum = 1; add to curSum dfs from all dimentsions and return curSum| 11.07 | Start time: 13:47 , End time: 14:04| Aided for 50% | Re-do
80 |Pacific Atlantic Water Flow| Create 2 new sets - pacific and atlantic. They will store visited coordinates (r & c). Run dfs on all rows at first, on 0th and on last column for both pacific and atlantic. Then run dfs on all cols, on 0th and last row, for both pacific and atlantic. 
Inside dfs add visited coordinates at first, then for each dir of dirs get new row and new col. Check if newRow & newCol are in bound(more than 0 and less than last, are not visited yet and heighs[newRow][newCol] >= heights[curRow][curCol] (backwards checking, since we're starting from edges). Finally, get intersect of pacific and atlantic sets | 13.07 | Start time: 11:11 , End time: 11:35| Aided| Re-do
81 |Surrounded Regions| run dfs on edges for each 0th and last row & 0th and last col, if its 'O' and mark it as 'S', to distinguish edge regions. Then run through each row and col again, and if it's O, set to X(is surrounded), if its S, set back to O (was on edge)| 14.07 | Start time: 11:45 , End time: 12:10| Aided for 50% | Re-do
82 |Rotting Oranges| use bfs, put already rottern into queue, and cunt how many fresh out there. Then explore queue and put into queue all ornages on same level. If no more fresh oranges return number of minutes | 15.07 | Start time: 12:31 , End time: | Aided for 50% | Re-check
83 |Course Schedule| Build adjacency list - for each value from numOfCourses, create empty array. then run dfs on each course (from 0 to numOfCourses) | 19.07 | Start time: 14:00, End time: 14:36| Aided | Re-do
84 |Course Schedule II| same as previous, push to path if dfs wasnt false | 19.07 | Start time: 15:00, End time: 16:00| Aided | Re-do
84 |Redundant Connection| Union Find, return if parent is same | 21.07 | Start time: 15:00, End time: 16:00| Aided | Re-do
84 |Word Ladder| BFS, inside it iterate through each character of the word and iterate through every character of alphabet and concatenate new words. If such word exist and isnt visited yet, push into queue | 22.07 | Start time: 15:00, End time: 16:00| Aided | Re-do
            for (var j = 0; j < word.length; j++) {
                for (var c = 97; c <= 122; c++) {
                    var newWord = word.slice(0, j) + String.fromCharCode(c) + word.slice(j + 1);
                    if (wordSet.has(newWord) && !visited.has(newWord)) {
                        queue.push(newWord);
                        visited.add(newWord);
                    }
                }
            }