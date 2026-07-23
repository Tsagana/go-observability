# BACKTRACKING
var solve = function(input) {
    const res = [];
    const path = [];

    function backtrack(start /* or other state */) {
        // 1. RECORD / BASE CASE
        if (/* goal reached */) {
            res.push([...path]);   // ALWAYS copy — path keeps mutating
            return;
        }
        // 2. PRUNE (optional but often essential)
        if (/* path can't lead to a solution */) return;

        // 3. CHOICES
        for (let i = start; i < input.length; i++) {
            // skip-duplicates line goes HERE when needed
            path.push(input[i]);          // choose
            backtrack(i + 1 /* or i */);  // explore
            path.pop();                   // un-choose
        }
    }

    backtrack(0);
    return res;
};

Backtracking on a list/array — inner for loop picks the next element, path is the state, i+1 / i / used controls reuse. Subsets, Combinations, Permutations.
Backtracking on a graph/grid — recursive calls are the choices (neighbors), position + goal-progress is the state, mark-and-restore the current node. Word Search, maze problems, Sudoku, N-Queens.
# DFS
// Recursive - most common in interviews
function dfsRecursive(root) {
  if (!root) return;

  // Pre-order: process node BEFORE children
  console.log(root.val);

  dfsRecursive(root.left);
  dfsRecursive(root.right);

  // Post-order: move console.log here (after children)
  // In-order: move console.log between left and right calls
}

// Iterative - use when recursion depth could be a problem
function dfsIterative(root) {
  if (!root) return;

  const stack = [root];

  while (stack.length) {
    const node = stack.pop();
    console.log(node.val);

    // Push right first so left is processed first (LIFO)
    if (node.right) stack.push(node.right);
    if (node.left) stack.push(node.left);
  }
}

# BFS
function bfs(root) {
  if (!root) return;

  const queue = [root];

  while (queue.length) {
    const node = queue.shift();
    console.log(node.val);

    if (node.left) queue.push(node.left);
    if (node.right) queue.push(node.right);
  }
}

// Level-by-level variant — very common in interviews
function bfsLevels(root) {
  if (!root) return [];

  const result = [];
  const queue = [root];

  while (queue.length) {
    const levelSize = queue.length;
    const level = [];

    for (let i = 0; i < levelSize; i++) {
      const node = queue.shift();
      level.push(node.val);

      if (node.left) queue.push(node.left);
      if (node.right) queue.push(node.right);
    }

    result.push(level);
  }

  return result; // [[1], [2, 3], [4, 5, 6, 7]]
}

# HEAPS
function topKFrequent(nums, k) {
  const freq = new Map();
  for (const n of nums) freq.set(n, (freq.get(n) || 0) + 1);

  // min-heap keyed by frequency
  const pq = new PriorityQueue((a, b) => a[1] - b[1]);

  for (const [num, count] of freq) {
    pq.push([num, count]);
    if (pq.size() > k) pq.pop(); // evict least frequent
  }

  return pq.heap.map(([num]) => num);
}

topKFrequent([1,1,1,2,2,3], 2); // [1, 2]

# INTERVALS
var intervalProblem = function(intervals) {
    // 1. SORT — almost always by start (sometimes by end, see below)
    intervals.sort((a, b) => a[0] - b[0]);

    // 2. INIT — track state across the sweep
    var res = [];          // or a counter, or a single value
    // optionally: var last = intervals[0]; res.push(last);

    // 3. SWEEP — one pass, compare current to previous state
    for (var i = 1; i < intervals.length; i++) {
        var cur = intervals[i];
        var prev = res[res.length - 1]; // or `last`, depending on problem

        if (prev[1] >= cur[0]) {
            // OVERLAP — handle based on problem:
            // merge:    prev[1] = Math.max(prev[1], cur[1]);
            // count:    overlaps++;
            // remove:   removed++; prev[1] = Math.min(prev[1], cur[1]); // keep the one ending earlier
            // can't attend / book: return false;
        } else {
            // NO OVERLAP — usually just record current
            res.push(cur);
        }
    }

    return res;
};

# BFS in grid
function bfs(grid, startRow, startCol) {
    const rows = grid.length;
    const cols = grid[0].length;
    const visited = new Set();
    const queue = [[startRow, startCol]];
    visited.add(`${startRow},${startCol}`);
    
    const directions = [[1,0], [-1,0], [0,1], [0,-1]];
    let steps = 0;
    
    while (queue.length > 0) {
        const size = queue.length; // process level by level
        for (let i = 0; i < size; i++) {
            const [row, col] = queue.shift();
            
            // process current cell here
            
            for (const [dr, dc] of directions) {
                const newRow = row + dr;
                const newCol = col + dc;
                const key = `${newRow},${newCol}`;
                
                if (
                    newRow >= 0 && newRow < rows &&
                    newCol >= 0 && newCol < cols &&
                    !visited.has(key) &&
                    grid[newRow][newCol] !== 'obstacle' // adjust condition
                ) {
                    visited.add(key);
                    queue.push([newRow, newCol]);
                }
            }
        }
        steps++;
    }
    
    return steps;
}

# DFS in grid
function dfs(grid, row, col, visited) {
    const rows = grid.length;
    const cols = grid[0].length;
    
    if (
        row < 0 || row >= rows ||
        col < 0 || col >= cols ||
        visited.has(`${row},${col}`) ||
        grid[row][col] !== 'valid' // adjust condition
    ) {
        return;
    }
    
    visited.add(`${row},${col}`);
    
    // process current cell here
    
    const directions = [[1,0], [-1,0], [0,1], [0,-1]];
    for (const [dr, dc] of directions) {
        dfs(grid, row + dr, col + dc, visited);
    }
}

# DFS in grid iterative
function dfsIterative(grid, startRow, startCol) {
    const rows = grid.length;
    const cols = grid[0].length;
    const visited = new Set();
    const stack = [[startRow, startCol]];
    
    const directions = [[1,0], [-1,0], [0,1], [0,-1]];
    
    while (stack.length > 0) {
        const [row, col] = stack.pop();
        const key = `${row},${col}`;
        
        if (visited.has(key)) continue;
        visited.add(key);
        
        // process current cell here
        
        for (const [dr, dc] of directions) {
            const newRow = row + dr;
            const newCol = col + dc;
            if (
                newRow >= 0 && newRow < rows &&
                newCol >= 0 && newCol < cols &&
                !visited.has(`${newRow},${newCol}`)
            ) {
                stack.push([newRow, newCol]);
            }
        }
    }
}

# BFS in graph
function bfs(graph, start) {
    const visited = new Set([start]);
    const queue = [start];
    let head = 0;
    
    while (head < queue.length) {
        const node = queue[head++];
        
        // process current node here
        
        for (const neighbor of graph[node]) {
            if (!visited.has(neighbor)) {
                visited.add(neighbor);
                queue.push(neighbor);
            }
        }
    }
}

# DFS in graph
function dfs(graph, node, visited = new Set()) {
    if (visited.has(node)) return;
    visited.add(node);
    
    // process current node here
    
    for (const neighbor of graph[node]) {
        dfs(graph, neighbor, visited);
    }
}

