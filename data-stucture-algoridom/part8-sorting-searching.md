# Part 8: Sorting and Searching Algorithms

## Table of Contents
1. [Introduction to Sorting](#introduction-sorting)
2. [Simple Sorting Algorithms](#simple-sorting)
3. [Efficient Sorting Algorithms](#efficient-sorting)
4. [Specialized Sorting Algorithms](#specialized-sorting)
5. [Searching Algorithms](#searching-algorithms)
6. [Advanced Searching Techniques](#advanced-searching)
7. [Sorting and Searching Applications](#applications)
8. [Practice Exercises](#practice-exercises)

---

## 1. Introduction to Sorting {#introduction-sorting}

### What is Sorting?
Sorting is the process of arranging elements in a specific order (ascending or descending) based on some criteria.

### Why is Sorting Important?
- **Data Organization**: Makes data easier to search and analyze
- **Algorithm Optimization**: Many algorithms work better on sorted data
- **Database Operations**: Indexing and query optimization
- **User Experience**: Presenting data in meaningful order

### Sorting Algorithm Properties
- **Stability**: Maintains relative order of equal elements
- **In-place**: Uses O(1) extra space
- **Adaptive**: Performs better on partially sorted data
- **Online**: Can sort data as it arrives

### Classification of Sorting Algorithms

#### By Time Complexity
- **O(n²)**: Bubble, Selection, Insertion Sort
- **O(n log n)**: Merge, Heap, Quick Sort (average)
- **O(n)**: Counting, Radix, Bucket Sort (under specific conditions)

#### By Space Complexity
- **In-place**: O(1) extra space
- **Out-of-place**: O(n) extra space

#### By Stability
- **Stable**: Merge, Insertion, Bubble Sort
- **Unstable**: Quick, Heap, Selection Sort

---

## 2. Simple Sorting Algorithms {#simple-sorting}

### Bubble Sort

```cpp
class BubbleSort {
public:
    // Basic bubble sort
    static void bubbleSort(vector<int>& arr) {
        int n = arr.size();
        
        for (int i = 0; i < n - 1; i++) {
            for (int j = 0; j < n - i - 1; j++) {
                if (arr[j] > arr[j + 1]) {
                    swap(arr[j], arr[j + 1]);
                }
            }
        }
    }
    
    // Optimized bubble sort (early termination)
    static void bubbleSortOptimized(vector<int>& arr) {
        int n = arr.size();
        bool swapped;
        
        for (int i = 0; i < n - 1; i++) {
            swapped = false;
            
            for (int j = 0; j < n - i - 1; j++) {
                if (arr[j] > arr[j + 1]) {
                    swap(arr[j], arr[j + 1]);
                    swapped = true;
                }
            }
            
            // If no swapping occurred, array is sorted
            if (!swapped) {
                break;
            }
        }
    }
    
    // Recursive bubble sort
    static void bubbleSortRecursive(vector<int>& arr, int n) {
        // Base case
        if (n == 1) {
            return;
        }
        
        // One pass of bubble sort
        for (int i = 0; i < n - 1; i++) {
            if (arr[i] > arr[i + 1]) {
                swap(arr[i], arr[i + 1]);
            }
        }
        
        // Recursively sort first n-1 elements
        bubbleSortRecursive(arr, n - 1);
    }
};
```

### Selection Sort

```cpp
class SelectionSort {
public:
    // Basic selection sort
    static void selectionSort(vector<int>& arr) {
        int n = arr.size();
        
        for (int i = 0; i < n - 1; i++) {
            int minIndex = i;
            
            // Find minimum element in remaining array
            for (int j = i + 1; j < n; j++) {
                if (arr[j] < arr[minIndex]) {
                    minIndex = j;
                }
            }
            
            // Swap minimum element with first element
            if (minIndex != i) {
                swap(arr[i], arr[minIndex]);
            }
        }
    }
    
    // Selection sort with custom comparator
    template<typename Compare>
    static void selectionSort(vector<int>& arr, Compare comp) {
        int n = arr.size();
        
        for (int i = 0; i < n - 1; i++) {
            int selectedIndex = i;
            
            for (int j = i + 1; j < n; j++) {
                if (comp(arr[j], arr[selectedIndex])) {
                    selectedIndex = j;
                }
            }
            
            if (selectedIndex != i) {
                swap(arr[i], arr[selectedIndex]);
            }
        }
    }
    
    // Bidirectional selection sort
    static void bidirectionalSelectionSort(vector<int>& arr) {
        int left = 0, right = arr.size() - 1;
        
        while (left < right) {
            int minIndex = left, maxIndex = left;
            
            // Find both minimum and maximum in current range
            for (int i = left; i <= right; i++) {
                if (arr[i] < arr[minIndex]) {
                    minIndex = i;
                }
                if (arr[i] > arr[maxIndex]) {
                    maxIndex = i;
                }
            }
            
            // Place minimum at left
            swap(arr[left], arr[minIndex]);
            
            // If maximum was at left position, it's now at minIndex
            if (maxIndex == left) {
                maxIndex = minIndex;
            }
            
            // Place maximum at right
            swap(arr[right], arr[maxIndex]);
            
            left++;
            right--;
        }
    }
};
```

### Insertion Sort

```cpp
class InsertionSort {
public:
    // Basic insertion sort
    static void insertionSort(vector<int>& arr) {
        int n = arr.size();
        
        for (int i = 1; i < n; i++) {
            int key = arr[i];
            int j = i - 1;
            
            // Move elements greater than key one position ahead
            while (j >= 0 && arr[j] > key) {
                arr[j + 1] = arr[j];
                j--;
            }
            
            arr[j + 1] = key;
        }
    }
    
    // Binary insertion sort
    static void binaryInsertionSort(vector<int>& arr) {
        int n = arr.size();
        
        for (int i = 1; i < n; i++) {
            int key = arr[i];
            int left = 0, right = i;
            
            // Binary search for insertion position
            while (left < right) {
                int mid = left + (right - left) / 2;
                if (arr[mid] > key) {
                    right = mid;
                } else {
                    left = mid + 1;
                }
            }
            
            // Shift elements and insert
            for (int j = i; j > left; j--) {
                arr[j] = arr[j - 1];
            }
            arr[left] = key;
        }
    }
    
    // Recursive insertion sort
    static void insertionSortRecursive(vector<int>& arr, int n) {
        // Base case
        if (n <= 1) {
            return;
        }
        
        // Sort first n-1 elements
        insertionSortRecursive(arr, n - 1);
        
        // Insert last element at correct position
        int last = arr[n - 1];
        int j = n - 2;
        
        while (j >= 0 && arr[j] > last) {
            arr[j + 1] = arr[j];
            j--;
        }
        
        arr[j + 1] = last;
    }
    
    // Insertion sort for linked list
    struct ListNode {
        int val;
        ListNode* next;
        ListNode(int x) : val(x), next(nullptr) {}
    };
    
    static ListNode* insertionSortList(ListNode* head) {
        if (!head || !head->next) {
            return head;
        }
        
        ListNode* dummy = new ListNode(0);
        ListNode* current = head;
        
        while (current) {
            ListNode* next = current->next;
            ListNode* prev = dummy;
            
            // Find insertion position
            while (prev->next && prev->next->val < current->val) {
                prev = prev->next;
            }
            
            // Insert current node
            current->next = prev->next;
            prev->next = current;
            
            current = next;
        }
        
        ListNode* result = dummy->next;
        delete dummy;
        return result;
    }
};
```

---

## 3. Efficient Sorting Algorithms {#efficient-sorting}

### Merge Sort

```cpp
class MergeSort {
public:
    // Basic merge sort
    static void mergeSort(vector<int>& arr, int left, int right) {
        if (left < right) {
            int mid = left + (right - left) / 2;
            
            mergeSort(arr, left, mid);
            mergeSort(arr, mid + 1, right);
            merge(arr, left, mid, right);
        }
    }
    
    // Wrapper function
    static void mergeSort(vector<int>& arr) {
        mergeSort(arr, 0, arr.size() - 1);
    }
    
    // Iterative merge sort
    static void mergeSortIterative(vector<int>& arr) {
        int n = arr.size();
        
        for (int size = 1; size < n; size *= 2) {
            for (int left = 0; left < n - 1; left += 2 * size) {
                int mid = min(left + size - 1, n - 1);
                int right = min(left + 2 * size - 1, n - 1);
                
                if (mid < right) {
                    merge(arr, left, mid, right);
                }
            }
        }
    }
    
    // In-place merge sort (using O(1) extra space)
    static void mergeSortInPlace(vector<int>& arr, int left, int right) {
        if (left < right) {
            int mid = left + (right - left) / 2;
            
            mergeSortInPlace(arr, left, mid);
            mergeSortInPlace(arr, mid + 1, right);
            mergeInPlace(arr, left, mid, right);
        }
    }
    
    // 3-way merge sort
    static void mergeSort3Way(vector<int>& arr, int left, int right) {
        if (left < right) {
            int mid1 = left + (right - left) / 3;
            int mid2 = left + 2 * (right - left) / 3;
            
            mergeSort3Way(arr, left, mid1);
            mergeSort3Way(arr, mid1 + 1, mid2);
            mergeSort3Way(arr, mid2 + 1, right);
            merge3Way(arr, left, mid1, mid2, right);
        }
    }
    
private:
    // Merge two sorted subarrays
    static void merge(vector<int>& arr, int left, int mid, int right) {
        vector<int> temp(right - left + 1);
        int i = left, j = mid + 1, k = 0;
        
        while (i <= mid && j <= right) {
            if (arr[i] <= arr[j]) {
                temp[k++] = arr[i++];
            } else {
                temp[k++] = arr[j++];
            }
        }
        
        while (i <= mid) {
            temp[k++] = arr[i++];
        }
        
        while (j <= right) {
            temp[k++] = arr[j++];
        }
        
        for (int i = 0; i < k; i++) {
            arr[left + i] = temp[i];
        }
    }
    
    // In-place merge (O(1) space but O(n²) time in worst case)
    static void mergeInPlace(vector<int>& arr, int left, int mid, int right) {
        int start2 = mid + 1;
        
        if (arr[mid] <= arr[start2]) {
            return; // Already sorted
        }
        
        while (left <= mid && start2 <= right) {
            if (arr[left] <= arr[start2]) {
                left++;
            } else {
                int value = arr[start2];
                int index = start2;
                
                // Shift elements
                while (index != left) {
                    arr[index] = arr[index - 1];
                    index--;
                }
                
                arr[left] = value;
                left++;
                mid++;
                start2++;
            }
        }
    }
    
    // 3-way merge
    static void merge3Way(vector<int>& arr, int left, int mid1, int mid2, int right) {
        vector<int> temp(right - left + 1);
        int i = left, j = mid1 + 1, k = mid2 + 1, l = 0;
        
        while (i <= mid1 && j <= mid2 && k <= right) {
            if (arr[i] <= arr[j] && arr[i] <= arr[k]) {
                temp[l++] = arr[i++];
            } else if (arr[j] <= arr[k]) {
                temp[l++] = arr[j++];
            } else {
                temp[l++] = arr[k++];
            }
        }
        
        while (i <= mid1 && j <= mid2) {
            if (arr[i] <= arr[j]) {
                temp[l++] = arr[i++];
            } else {
                temp[l++] = arr[j++];
            }
        }
        
        while (i <= mid1 && k <= right) {
            if (arr[i] <= arr[k]) {
                temp[l++] = arr[i++];
            } else {
                temp[l++] = arr[k++];
            }
        }
        
        while (j <= mid2 && k <= right) {
            if (arr[j] <= arr[k]) {
                temp[l++] = arr[j++];
            } else {
                temp[l++] = arr[k++];
            }
        }
        
        while (i <= mid1) temp[l++] = arr[i++];
        while (j <= mid2) temp[l++] = arr[j++];
        while (k <= right) temp[l++] = arr[k++];
        
        for (int i = 0; i < l; i++) {
            arr[left + i] = temp[i];
        }
    }
};
```

### Quick Sort

```cpp
class QuickSort {
public:
    // Basic quick sort
    static void quickSort(vector<int>& arr, int low, int high) {
        if (low < high) {
            int pi = partition(arr, low, high);
            
            quickSort(arr, low, pi - 1);
            quickSort(arr, pi + 1, high);
        }
    }
    
    // Wrapper function
    static void quickSort(vector<int>& arr) {
        quickSort(arr, 0, arr.size() - 1);
    }
    
    // Iterative quick sort
    static void quickSortIterative(vector<int>& arr) {
        int n = arr.size();
        stack<pair<int, int>> stk;
        
        stk.push({0, n - 1});
        
        while (!stk.empty()) {
            auto [low, high] = stk.top();
            stk.pop();
            
            if (low < high) {
                int pi = partition(arr, low, high);
                
                stk.push({low, pi - 1});
                stk.push({pi + 1, high});
            }
        }
    }
    
    // 3-way quick sort (handles duplicates efficiently)
    static void quickSort3Way(vector<int>& arr, int low, int high) {
        if (low < high) {
            auto [lt, gt] = partition3Way(arr, low, high);
            
            quickSort3Way(arr, low, lt - 1);
            quickSort3Way(arr, gt + 1, high);
        }
    }
    
    // Randomized quick sort
    static void quickSortRandomized(vector<int>& arr, int low, int high) {
        if (low < high) {
            // Random pivot selection
            int randomIndex = low + rand() % (high - low + 1);
            swap(arr[randomIndex], arr[high]);
            
            int pi = partition(arr, low, high);
            
            quickSortRandomized(arr, low, pi - 1);
            quickSortRandomized(arr, pi + 1, high);
        }
    }
    
    // Hybrid quick sort (switches to insertion sort for small arrays)
    static void quickSortHybrid(vector<int>& arr, int low, int high) {
        const int THRESHOLD = 10;
        
        if (high - low + 1 <= THRESHOLD) {
            insertionSortRange(arr, low, high);
        } else if (low < high) {
            int pi = partition(arr, low, high);
            
            quickSortHybrid(arr, low, pi - 1);
            quickSortHybrid(arr, pi + 1, high);
        }
    }
    
private:
    // Lomuto partition scheme
    static int partition(vector<int>& arr, int low, int high) {
        int pivot = arr[high];
        int i = low - 1;
        
        for (int j = low; j < high; j++) {
            if (arr[j] <= pivot) {
                i++;
                swap(arr[i], arr[j]);
            }
        }
        
        swap(arr[i + 1], arr[high]);
        return i + 1;
    }
    
    // Hoare partition scheme
    static int partitionHoare(vector<int>& arr, int low, int high) {
        int pivot = arr[low];
        int i = low - 1, j = high + 1;
        
        while (true) {
            do {
                i++;
            } while (arr[i] < pivot);
            
            do {
                j--;
            } while (arr[j] > pivot);
            
            if (i >= j) {
                return j;
            }
            
            swap(arr[i], arr[j]);
        }
    }
    
    // 3-way partition (Dutch National Flag)
    static pair<int, int> partition3Way(vector<int>& arr, int low, int high) {
        int pivot = arr[low];
        int lt = low, gt = high, i = low;
        
        while (i <= gt) {
            if (arr[i] < pivot) {
                swap(arr[lt++], arr[i++]);
            } else if (arr[i] > pivot) {
                swap(arr[i], arr[gt--]);
            } else {
                i++;
            }
        }
        
        return {lt, gt};
    }
    
    // Insertion sort for small ranges
    static void insertionSortRange(vector<int>& arr, int low, int high) {
        for (int i = low + 1; i <= high; i++) {
            int key = arr[i];
            int j = i - 1;
            
            while (j >= low && arr[j] > key) {
                arr[j + 1] = arr[j];
                j--;
            }
            
            arr[j + 1] = key;
        }
    }
};
```

### Heap Sort

```cpp
class HeapSort {
public:
    // Basic heap sort
    static void heapSort(vector<int>& arr) {
        int n = arr.size();
        
        // Build max heap
        for (int i = n / 2 - 1; i >= 0; i--) {
            heapify(arr, n, i);
        }
        
        // Extract elements from heap one by one
        for (int i = n - 1; i > 0; i--) {
            swap(arr[0], arr[i]);
            heapify(arr, i, 0);
        }
    }
    
    // Min heap sort (descending order)
    static void heapSortDescending(vector<int>& arr) {
        int n = arr.size();
        
        // Build min heap
        for (int i = n / 2 - 1; i >= 0; i--) {
            minHeapify(arr, n, i);
        }
        
        // Extract elements from heap one by one
        for (int i = n - 1; i > 0; i--) {
            swap(arr[0], arr[i]);
            minHeapify(arr, i, 0);
        }
    }
    
    // Iterative heap sort
    static void heapSortIterative(vector<int>& arr) {
        int n = arr.size();
        
        // Build max heap iteratively
        for (int i = 1; i < n; i++) {
            int child = i;
            while (child > 0) {
                int parent = (child - 1) / 2;
                if (arr[child] > arr[parent]) {
                    swap(arr[child], arr[parent]);
                    child = parent;
                } else {
                    break;
                }
            }
        }
        
        // Extract elements
        for (int i = n - 1; i > 0; i--) {
            swap(arr[0], arr[i]);
            heapifyIterative(arr, i, 0);
        }
    }
    
private:
    // Max heapify
    static void heapify(vector<int>& arr, int n, int i) {
        int largest = i;
        int left = 2 * i + 1;
        int right = 2 * i + 2;
        
        if (left < n && arr[left] > arr[largest]) {
            largest = left;
        }
        
        if (right < n && arr[right] > arr[largest]) {
            largest = right;
        }
        
        if (largest != i) {
            swap(arr[i], arr[largest]);
            heapify(arr, n, largest);
        }
    }
    
    // Min heapify
    static void minHeapify(vector<int>& arr, int n, int i) {
        int smallest = i;
        int left = 2 * i + 1;
        int right = 2 * i + 2;
        
        if (left < n && arr[left] < arr[smallest]) {
            smallest = left;
        }
        
        if (right < n && arr[right] < arr[smallest]) {
            smallest = right;
        }
        
        if (smallest != i) {
            swap(arr[i], arr[smallest]);
            minHeapify(arr, n, smallest);
        }
    }
    
    // Iterative heapify
    static void heapifyIterative(vector<int>& arr, int n, int i) {
        while (true) {
            int largest = i;
            int left = 2 * i + 1;
            int right = 2 * i + 2;
            
            if (left < n && arr[left] > arr[largest]) {
                largest = left;
            }
            
            if (right < n && arr[right] > arr[largest]) {
                largest = right;
            }
            
            if (largest != i) {
                swap(arr[i], arr[largest]);
                i = largest;
            } else {
                break;
            }
        }
    }
};
```

---

## 4. Specialized Sorting Algorithms {#specialized-sorting}

### Counting Sort

```cpp
class CountingSort {
public:
    // Basic counting sort (for non-negative integers)
    static void countingSort(vector<int>& arr) {
        if (arr.empty()) return;
        
        int maxVal = *max_element(arr.begin(), arr.end());
        int minVal = *min_element(arr.begin(), arr.end());
        int range = maxVal - minVal + 1;
        
        vector<int> count(range, 0);
        vector<int> output(arr.size());
        
        // Count occurrences
        for (int num : arr) {
            count[num - minVal]++;
        }
        
        // Calculate cumulative count
        for (int i = 1; i < range; i++) {
            count[i] += count[i - 1];
        }
        
        // Build output array
        for (int i = arr.size() - 1; i >= 0; i--) {
            output[count[arr[i] - minVal] - 1] = arr[i];
            count[arr[i] - minVal]--;
        }
        
        // Copy back to original array
        for (int i = 0; i < arr.size(); i++) {
            arr[i] = output[i];
        }
    }
    
    // Counting sort for characters
    static void countingSortChars(string& str) {
        vector<int> count(256, 0);
        
        // Count characters
        for (char c : str) {
            count[c]++;
        }
        
        // Rebuild string
        int index = 0;
        for (int i = 0; i < 256; i++) {
            while (count[i]-- > 0) {
                str[index++] = (char)i;
            }
        }
    }
    
    // Counting sort with custom key function
    template<typename T, typename KeyFunc>
    static void countingSortWithKey(vector<T>& arr, KeyFunc getKey, int maxKey) {
        vector<int> count(maxKey + 1, 0);
        vector<T> output(arr.size());
        
        // Count occurrences
        for (const T& item : arr) {
            count[getKey(item)]++;
        }
        
        // Calculate cumulative count
        for (int i = 1; i <= maxKey; i++) {
            count[i] += count[i - 1];
        }
        
        // Build output array
        for (int i = arr.size() - 1; i >= 0; i--) {
            int key = getKey(arr[i]);
            output[count[key] - 1] = arr[i];
            count[key]--;
        }
        
        // Copy back
        arr = output;
    }
};
```

### Radix Sort

```cpp
class RadixSort {
public:
    // Radix sort for non-negative integers
    static void radixSort(vector<int>& arr) {
        if (arr.empty()) return;
        
        int maxVal = *max_element(arr.begin(), arr.end());
        
        // Do counting sort for every digit
        for (int exp = 1; maxVal / exp > 0; exp *= 10) {
            countingSortByDigit(arr, exp);
        }
    }
    
    // Radix sort for strings
    static void radixSortStrings(vector<string>& arr) {
        if (arr.empty()) return;
        
        int maxLen = 0;
        for (const string& str : arr) {
            maxLen = max(maxLen, (int)str.length());
        }
        
        // Sort by each character position from right to left
        for (int pos = maxLen - 1; pos >= 0; pos--) {
            countingSortByChar(arr, pos);
        }
    }
    
    // LSD (Least Significant Digit) radix sort
    static void radixSortLSD(vector<int>& arr, int digits) {
        for (int i = 0; i < digits; i++) {
            countingSortByDigit(arr, pow(10, i));
        }
    }
    
    // MSD (Most Significant Digit) radix sort
    static void radixSortMSD(vector<int>& arr, int left, int right, int exp) {
        if (left >= right || exp <= 0) return;
        
        vector<vector<int>> buckets(10);
        
        // Distribute into buckets
        for (int i = left; i <= right; i++) {
            int digit = (arr[i] / exp) % 10;
            buckets[digit].push_back(arr[i]);
        }
        
        // Copy back and recursively sort each bucket
        int index = left;
        for (int i = 0; i < 10; i++) {
            int bucketStart = index;
            for (int num : buckets[i]) {
                arr[index++] = num;
            }
            
            if (buckets[i].size() > 1) {
                radixSortMSD(arr, bucketStart, index - 1, exp / 10);
            }
        }
    }
    
private:
    // Counting sort by specific digit
    static void countingSortByDigit(vector<int>& arr, int exp) {
        vector<int> output(arr.size());
        vector<int> count(10, 0);
        
        // Count occurrences of each digit
        for (int num : arr) {
            count[(num / exp) % 10]++;
        }
        
        // Calculate cumulative count
        for (int i = 1; i < 10; i++) {
            count[i] += count[i - 1];
        }
        
        // Build output array
        for (int i = arr.size() - 1; i >= 0; i--) {
            int digit = (arr[i] / exp) % 10;
            output[count[digit] - 1] = arr[i];
            count[digit]--;
        }
        
        // Copy back
        for (int i = 0; i < arr.size(); i++) {
            arr[i] = output[i];
        }
    }
    
    // Counting sort by character at specific position
    static void countingSortByChar(vector<string>& arr, int pos) {
        vector<string> output(arr.size());
        vector<int> count(256, 0);
        
        // Count characters (treat missing chars as '\0')
        for (const string& str : arr) {
            char ch = (pos < str.length()) ? str[pos] : '\0';
            count[ch]++;
        }
        
        // Calculate cumulative count
        for (int i = 1; i < 256; i++) {
            count[i] += count[i - 1];
        }
        
        // Build output array
        for (int i = arr.size() - 1; i >= 0; i--) {
            char ch = (pos < arr[i].length()) ? arr[i][pos] : '\0';
            output[count[ch] - 1] = arr[i];
            count[ch]--;
        }
        
        // Copy back
        arr = output;
    }
};
```

### Bucket Sort

```cpp
class BucketSort {
public:
    // Bucket sort for floating point numbers in [0, 1)
    static void bucketSort(vector<double>& arr) {
        if (arr.empty()) return;
        
        int n = arr.size();
        vector<vector<double>> buckets(n);
        
        // Put array elements in different buckets
        for (double num : arr) {
            int bucketIndex = n * num;
            buckets[bucketIndex].push_back(num);
        }
        
        // Sort individual buckets
        for (auto& bucket : buckets) {
            sort(bucket.begin(), bucket.end());
        }
        
        // Concatenate all buckets
        int index = 0;
        for (const auto& bucket : buckets) {
            for (double num : bucket) {
                arr[index++] = num;
            }
        }
    }
    
    // Bucket sort for integers
    static void bucketSortIntegers(vector<int>& arr) {
        if (arr.empty()) return;
        
        int minVal = *min_element(arr.begin(), arr.end());
        int maxVal = *max_element(arr.begin(), arr.end());
        int range = maxVal - minVal + 1;
        int bucketCount = sqrt(arr.size());
        
        if (bucketCount == 0) bucketCount = 1;
        
        vector<vector<int>> buckets(bucketCount);
        
        // Distribute elements into buckets
        for (int num : arr) {
            int bucketIndex = (num - minVal) * bucketCount / range;
            if (bucketIndex >= bucketCount) bucketIndex = bucketCount - 1;
            buckets[bucketIndex].push_back(num);
        }
        
        // Sort individual buckets
        for (auto& bucket : buckets) {
            sort(bucket.begin(), bucket.end());
        }
        
        // Concatenate all buckets
        int index = 0;
        for (const auto& bucket : buckets) {
            for (int num : bucket) {
                arr[index++] = num;
            }
        }
    }
    
    // Generic bucket sort with custom bucket function
    template<typename T, typename BucketFunc>
    static void bucketSortGeneric(vector<T>& arr, int bucketCount, BucketFunc getBucket) {
        if (arr.empty()) return;
        
        vector<vector<T>> buckets(bucketCount);
        
        // Distribute elements
        for (const T& item : arr) {
            int bucketIndex = getBucket(item);
            buckets[bucketIndex].push_back(item);
        }
        
        // Sort individual buckets
        for (auto& bucket : buckets) {
            sort(bucket.begin(), bucket.end());
        }
        
        // Concatenate
        int index = 0;
        for (const auto& bucket : buckets) {
            for (const T& item : bucket) {
                arr[index++] = item;
            }
        }
    }
};
```

---

## 5. Searching Algorithms {#searching-algorithms}

### Linear Search

```cpp
class LinearSearch {
public:
    // Basic linear search
    static int linearSearch(const vector<int>& arr, int target) {
        for (int i = 0; i < arr.size(); i++) {
            if (arr[i] == target) {
                return i;
            }
        }
        return -1;
    }
    
    // Linear search with custom comparator
    template<typename T, typename Compare>
    static int linearSearch(const vector<T>& arr, const T& target, Compare comp) {
        for (int i = 0; i < arr.size(); i++) {
            if (comp(arr[i], target)) {
                return i;
            }
        }
        return -1;
    }
    
    // Find all occurrences
    static vector<int> linearSearchAll(const vector<int>& arr, int target) {
        vector<int> indices;
        for (int i = 0; i < arr.size(); i++) {
            if (arr[i] == target) {
                indices.push_back(i);
            }
        }
        return indices;
    }
    
    // Linear search from both ends
    static int linearSearchBidirectional(const vector<int>& arr, int target) {
        int left = 0, right = arr.size() - 1;
        
        while (left <= right) {
            if (arr[left] == target) {
                return left;
            }
            if (arr[right] == target) {
                return right;
            }
            left++;
            right--;
        }
        
        return -1;
    }
    
    // Recursive linear search
    static int linearSearchRecursive(const vector<int>& arr, int target, int index = 0) {
        if (index >= arr.size()) {
            return -1;
        }
        
        if (arr[index] == target) {
            return index;
        }
        
        return linearSearchRecursive(arr, target, index + 1);
    }
};
```

### Binary Search

```cpp
class BinarySearch {
public:
    // Basic binary search (iterative)
    static int binarySearch(const vector<int>& arr, int target) {
        int left = 0, right = arr.size() - 1;
        
        while (left <= right) {
            int mid = left + (right - left) / 2;
            
            if (arr[mid] == target) {
                return mid;
            } else if (arr[mid] < target) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
        
        return -1;
    }
    
    // Recursive binary search
    static int binarySearchRecursive(const vector<int>& arr, int target, int left, int right) {
        if (left > right) {
            return -1;
        }
        
        int mid = left + (right - left) / 2;
        
        if (arr[mid] == target) {
            return mid;
        } else if (arr[mid] < target) {
            return binarySearchRecursive(arr, target, mid + 1, right);
        } else {
            return binarySearchRecursive(arr, target, left, mid - 1);
        }
    }
    
    // Find first occurrence
    static int findFirst(const vector<int>& arr, int target) {
        int left = 0, right = arr.size() - 1;
        int result = -1;
        
        while (left <= right) {
            int mid = left + (right - left) / 2;
            
            if (arr[mid] == target) {
                result = mid;
                right = mid - 1; // Continue searching left
            } else if (arr[mid] < target) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
        
        return result;
    }
    
    // Find last occurrence
    static int findLast(const vector<int>& arr, int target) {
        int left = 0, right = arr.size() - 1;
        int result = -1;
        
        while (left <= right) {
            int mid = left + (right - left) / 2;
            
            if (arr[mid] == target) {
                result = mid;
                left = mid + 1; // Continue searching right
            } else if (arr[mid] < target) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
        
        return result;
    }
    
    // Count occurrences
    static int countOccurrences(const vector<int>& arr, int target) {
        int first = findFirst(arr, target);
        if (first == -1) return 0;
        
        int last = findLast(arr, target);
        return last - first + 1;
    }
    
    // Find insertion position (lower_bound)
    static int lowerBound(const vector<int>& arr, int target) {
        int left = 0, right = arr.size();
        
        while (left < right) {
            int mid = left + (right - left) / 2;
            
            if (arr[mid] < target) {
                left = mid + 1;
            } else {
                right = mid;
            }
        }
        
        return left;
    }
    
    // Find upper bound
    static int upperBound(const vector<int>& arr, int target) {
        int left = 0, right = arr.size();
        
        while (left < right) {
            int mid = left + (right - left) / 2;
            
            if (arr[mid] <= target) {
                left = mid + 1;
            } else {
                right = mid;
            }
        }
        
        return left;
    }
    
    // Binary search on answer (find minimum value that satisfies condition)
    template<typename Predicate>
    static int binarySearchOnAnswer(int left, int right, Predicate canAchieve) {
        int result = right + 1;
        
        while (left <= right) {
            int mid = left + (right - left) / 2;
            
            if (canAchieve(mid)) {
                result = mid;
                right = mid - 1;
            } else {
                left = mid + 1;
            }
        }
        
        return result;
    }
};
```

---

## 6. Advanced Searching Techniques {#advanced-searching}

### Ternary Search

```cpp
class TernarySearch {
public:
    // Ternary search for unimodal functions
    static double ternarySearch(function<double(double)> f, double left, double right, double eps = 1e-9) {
        while (right - left > eps) {
            double m1 = left + (right - left) / 3;
            double m2 = right - (right - left) / 3;
            
            if (f(m1) < f(m2)) {
                left = m1;
            } else {
                right = m2;
            }
        }
        
        return (left + right) / 2;
    }
    
    // Ternary search on discrete array
    static int ternarySearchDiscrete(const vector<int>& arr, int target) {
        int left = 0, right = arr.size() - 1;
        
        while (left <= right) {
            if (right - left < 2) {
                if (arr[left] == target) return left;
                if (arr[right] == target) return right;
                return -1;
            }
            
            int m1 = left + (right - left) / 3;
            int m2 = right - (right - left) / 3;
            
            if (arr[m1] == target) return m1;
            if (arr[m2] == target) return m2;
            
            if (target < arr[m1]) {
                right = m1 - 1;
            } else if (target > arr[m2]) {
                left = m2 + 1;
            } else {
                left = m1 + 1;
                right = m2 - 1;
            }
        }
        
        return -1;
    }
};
```

### Exponential Search

```cpp
class ExponentialSearch {
public:
    // Exponential search (for unbounded arrays)
    static int exponentialSearch(const vector<int>& arr, int target) {
        if (arr.empty()) return -1;
        if (arr[0] == target) return 0;
        
        // Find range for binary search
        int bound = 1;
        while (bound < arr.size() && arr[bound] < target) {
            bound *= 2;
        }
        
        // Binary search in found range
        int left = bound / 2;
        int right = min(bound, (int)arr.size() - 1);
        
        return BinarySearch::binarySearchRecursive(arr, target, left, right);
    }
    
    // Exponential search with custom step
    static int exponentialSearchCustom(const vector<int>& arr, int target, int step = 2) {
        if (arr.empty()) return -1;
        if (arr[0] == target) return 0;
        
        int bound = 1;
        while (bound < arr.size() && arr[bound] < target) {
            bound *= step;
        }
        
        int left = bound / step;
        int right = min(bound, (int)arr.size() - 1);
        
        return BinarySearch::binarySearchRecursive(arr, target, left, right);
    }
};
```

### Interpolation Search

```cpp
class InterpolationSearch {
public:
    // Interpolation search (for uniformly distributed data)
    static int interpolationSearch(const vector<int>& arr, int target) {
        int left = 0, right = arr.size() - 1;
        
        while (left <= right && target >= arr[left] && target <= arr[right]) {
            if (left == right) {
                return (arr[left] == target) ? left : -1;
            }
            
            // Interpolation formula
            int pos = left + ((double)(target - arr[left]) / (arr[right] - arr[left])) * (right - left);
            
            if (arr[pos] == target) {
                return pos;
            } else if (arr[pos] < target) {
                left = pos + 1;
            } else {
                right = pos - 1;
            }
        }
        
        return -1;
    }
    
    // Interpolation search with bounds checking
    static int interpolationSearchSafe(const vector<int>& arr, int target) {
        int left = 0, right = arr.size() - 1;
        
        while (left <= right && target >= arr[left] && target <= arr[right]) {
            if (left == right) {
                return (arr[left] == target) ? left : -1;
            }
            
            if (arr[right] == arr[left]) {
                // All elements in range are equal
                return (arr[left] == target) ? left : -1;
            }
            
            // Calculate position with bounds checking
            double ratio = (double)(target - arr[left]) / (arr[right] - arr[left]);
            int pos = left + ratio * (right - left);
            
            // Ensure pos is within bounds
            pos = max(left, min(right, pos));
            
            if (arr[pos] == target) {
                return pos;
            } else if (arr[pos] < target) {
                left = pos + 1;
            } else {
                right = pos - 1;
            }
        }
        
        return -1;
    }
};
```

### Jump Search

```cpp
class JumpSearch {
public:
    // Basic jump search
    static int jumpSearch(const vector<int>& arr, int target) {
        int n = arr.size();
        int step = sqrt(n);
        int prev = 0;
        
        // Find block where element may be present
        while (arr[min(step, n) - 1] < target) {
            prev = step;
            step += sqrt(n);
            if (prev >= n) {
                return -1;
            }
        }
        
        // Linear search in identified block
        while (arr[prev] < target) {
            prev++;
            if (prev == min(step, n)) {
                return -1;
            }
        }
        
        return (arr[prev] == target) ? prev : -1;
    }
    
    // Jump search with custom step size
    static int jumpSearchCustom(const vector<int>& arr, int target, int jumpSize) {
        int n = arr.size();
        int step = jumpSize;
        int prev = 0;
        
        while (arr[min(step, n) - 1] < target) {
            prev = step;
            step += jumpSize;
            if (prev >= n) {
                return -1;
            }
        }
        
        while (arr[prev] < target) {
            prev++;
            if (prev == min(step, n)) {
                return -1;
            }
        }
        
        return (arr[prev] == target) ? prev : -1;
    }
    
    // Optimal jump search (finds optimal jump size)
    static int jumpSearchOptimal(const vector<int>& arr, int target) {
        int n = arr.size();
        int optimalJump = sqrt(n);
        
        return jumpSearchCustom(arr, target, optimalJump);
    }
};
```

---

## 7. Sorting and Searching Applications {#applications}

### Practical Applications

```cpp
class SortingSearchingApplications {
public:
    // Find kth largest element
    static int findKthLargest(vector<int>& nums, int k) {
        // Using quickselect algorithm
        return quickSelect(nums, 0, nums.size() - 1, nums.size() - k);
    }
    
    // Find median of two sorted arrays
    static double findMedianSortedArrays(vector<int>& nums1, vector<int>& nums2) {
        if (nums1.size() > nums2.size()) {
            return findMedianSortedArrays(nums2, nums1);
        }
        
        int m = nums1.size(), n = nums2.size();
        int left = 0, right = m;
        
        while (left <= right) {
            int partitionX = (left + right) / 2;
            int partitionY = (m + n + 1) / 2 - partitionX;
            
            int maxLeftX = (partitionX == 0) ? INT_MIN : nums1[partitionX - 1];
            int minRightX = (partitionX == m) ? INT_MAX : nums1[partitionX];
            
            int maxLeftY = (partitionY == 0) ? INT_MIN : nums2[partitionY - 1];
            int minRightY = (partitionY == n) ? INT_MAX : nums2[partitionY];
            
            if (maxLeftX <= minRightY && maxLeftY <= minRightX) {
                if ((m + n) % 2 == 0) {
                    return (max(maxLeftX, maxLeftY) + min(minRightX, minRightY)) / 2.0;
                } else {
                    return max(maxLeftX, maxLeftY);
                }
            } else if (maxLeftX > minRightY) {
                right = partitionX - 1;
            } else {
                left = partitionX + 1;
            }
        }
        
        return 0.0;
    }
    
    // Search in rotated sorted array
    static int searchRotated(vector<int>& nums, int target) {
        int left = 0, right = nums.size() - 1;
        
        while (left <= right) {
            int mid = left + (right - left) / 2;
            
            if (nums[mid] == target) {
                return mid;
            }
            
            // Left half is sorted
            if (nums[left] <= nums[mid]) {
                if (target >= nums[left] && target < nums[mid]) {
                    right = mid - 1;
                } else {
                    left = mid + 1;
                }
            }
            // Right half is sorted
            else {
                if (target > nums[mid] && target <= nums[right]) {
                    left = mid + 1;
                } else {
                    right = mid - 1;
                }
            }
        }
        
        return -1;
    }
    
    // Find peak element
    static int findPeakElement(vector<int>& nums) {
        int left = 0, right = nums.size() - 1;
        
        while (left < right) {
            int mid = left + (right - left) / 2;
            
            if (nums[mid] > nums[mid + 1]) {
                right = mid;
            } else {
                left = mid + 1;
            }
        }
        
        return left;
    }
    
    // Merge k sorted arrays
    static vector<int> mergeKSortedArrays(vector<vector<int>>& arrays) {
        priority_queue<tuple<int, int, int>, vector<tuple<int, int, int>>, greater<tuple<int, int, int>>> pq;
        
        // Initialize heap with first element of each array
        for (int i = 0; i < arrays.size(); i++) {
            if (!arrays[i].empty()) {
                pq.push({arrays[i][0], i, 0});
            }
        }
        
        vector<int> result;
        
        while (!pq.empty()) {
            auto [val, arrayIdx, elemIdx] = pq.top();
            pq.pop();
            
            result.push_back(val);
            
            // Add next element from same array
            if (elemIdx + 1 < arrays[arrayIdx].size()) {
                pq.push({arrays[arrayIdx][elemIdx + 1], arrayIdx, elemIdx + 1});
            }
        }
        
        return result;
    }
    
private:
    // QuickSelect for finding kth element
    static int quickSelect(vector<int>& nums, int left, int right, int k) {
        if (left == right) {
            return nums[left];
        }
        
        int pivotIndex = partition(nums, left, right);
        
        if (k == pivotIndex) {
            return nums[k];
        } else if (k < pivotIndex) {
            return quickSelect(nums, left, pivotIndex - 1, k);
        } else {
            return quickSelect(nums, pivotIndex + 1, right, k);
        }
    }
    
    static int partition(vector<int>& nums, int left, int right) {
        int pivot = nums[right];
        int i = left;
        
        for (int j = left; j < right; j++) {
            if (nums[j] <= pivot) {
                swap(nums[i], nums[j]);
                i++;
            }
        }
        
        swap(nums[i], nums[right]);
        return i;
    }
};
```

---

## 8. Practice Exercises {#practice-exercises}

### Basic Sorting Exercises
1. **Implementation Practice**:
   - Implement all basic sorting algorithms
   - Compare performance on different data sets
   - Implement stable versions of unstable algorithms

2. **Sorting Variations**:
   - Sort array of strings by length
   - Sort 2D points by distance from origin
   - Sort intervals by start time
   - Custom comparator implementations

3. **Optimization Challenges**:
   - Optimize bubble sort with early termination
   - Implement hybrid sorting algorithms
   - Memory-efficient sorting for large datasets

### Intermediate Sorting Exercises
1. **Specialized Sorting**:
   - Sort array with only 0s, 1s, and 2s (Dutch National Flag)
   - Sort array by frequency of elements
   - Sort nearly sorted array efficiently
   - External sorting for files larger than memory

2. **Advanced Problems**:
   - Merge k sorted linked lists
   - Sort array with limited swaps
   - Pancake sorting problem
   - Wiggle sort (arrange in wave form)

### Basic Searching Exercises
1. **Search Implementations**:
   - Implement all search algorithms
   - Search in 2D matrix
   - Find element in infinite sorted array
   - Search for range in sorted array

2. **Binary Search Variations**:
   - Find square root using binary search
   - Search in rotated sorted array
   - Find minimum in rotated sorted array
   - Search for target in mountain array

### Advanced Searching Exercises
1. **Complex Search Problems**:
   - Find median in data stream
   - Kth largest element in stream
   - Find duplicate number in array
   - Search in 2D sorted matrix

2. **Optimization Problems**:
   - Minimize maximum distance (binary search on answer)
   - Allocate minimum pages (binary search application)
   - Find minimum time to complete tasks
   - Capacity to ship packages within D days

### Complete Example Program

```cpp
#include <iostream>
#include <vector>
#include <algorithm>
#include <chrono>
#include <random>
using namespace std;
using namespace chrono;

class SortingBenchmark {
public:
    static void benchmarkSortingAlgorithms() {
        vector<int> sizes = {1000, 5000, 10000, 50000};
        
        for (int size : sizes) {
            cout << "\n=== Benchmarking with " << size << " elements ===" << endl;
            
            // Generate random data
            vector<int> data = generateRandomData(size);
            
            // Test different sorting algorithms
            testSortingAlgorithm("Bubble Sort", data, [](vector<int>& arr) {
                BubbleSort::bubbleSortOptimized(arr);
            });
            
            testSortingAlgorithm("Selection Sort", data, [](vector<int>& arr) {
                SelectionSort::selectionSort(arr);
            });
            
            testSortingAlgorithm("Insertion Sort", data, [](vector<int>& arr) {
                InsertionSort::insertionSort(arr);
            });
            
            testSortingAlgorithm("Merge Sort", data, [](vector<int>& arr) {
                MergeSort::mergeSort(arr);
            });
            
            testSortingAlgorithm("Quick Sort", data, [](vector<int>& arr) {
                QuickSort::quickSort(arr);
            });
            
            testSortingAlgorithm("Heap Sort", data, [](vector<int>& arr) {
                HeapSort::heapSort(arr);
            });
            
            testSortingAlgorithm("STL Sort", data, [](vector<int>& arr) {
                sort(arr.begin(), arr.end());
            });
        }
    }
    
    static void benchmarkSearchingAlgorithms() {
        vector<int> sizes = {1000, 10000, 100000, 1000000};
        
        for (int size : sizes) {
            cout << "\n=== Search Benchmarking with " << size << " elements ===" << endl;
            
            // Generate sorted data
            vector<int> data = generateSortedData(size);
            int target = data[size / 2]; // Middle element
            
            // Test different search algorithms
            testSearchAlgorithm("Linear Search", data, target, [](const vector<int>& arr, int t) {
                return LinearSearch::linearSearch(arr, t);
            });
            
            testSearchAlgorithm("Binary Search", data, target, [](const vector<int>& arr, int t) {
                return BinarySearch::binarySearch(arr, t);
            });
            
            testSearchAlgorithm("Jump Search", data, target, [](const vector<int>& arr, int t) {
                return JumpSearch::jumpSearch(arr, t);
            });
            
            testSearchAlgorithm("Interpolation Search", data, target, [](const vector<int>& arr, int t) {
                return InterpolationSearch::interpolationSearch(arr, t);
            });
        }
    }
    
private:
    static vector<int> generateRandomData(int size) {
        vector<int> data(size);
        random_device rd;
        mt19937 gen(rd());
        uniform_int_distribution<> dis(1, size * 10);
        
        for (int& num : data) {
            num = dis(gen);
        }
        
        return data;
    }
    
    static vector<int> generateSortedData(int size) {
        vector<int> data(size);
        for (int i = 0; i < size; i++) {
            data[i] = i * 2; // Even numbers
        }
        return data;
    }
    
    template<typename SortFunc>
    static void testSortingAlgorithm(const string& name, vector<int> data, SortFunc sortFunc) {
        auto start = high_resolution_clock::now();
        sortFunc(data);
        auto end = high_resolution_clock::now();
        
        auto duration = duration_cast<microseconds>(end - start);
        cout << name << ": " << duration.count() << " microseconds" << endl;
        
        // Verify sorting
        if (!is_sorted(data.begin(), data.end())) {
            cout << "ERROR: " << name << " did not sort correctly!" << endl;
        }
    }
    
    template<typename SearchFunc>
    static void testSearchAlgorithm(const string& name, const vector<int>& data, int target, SearchFunc searchFunc) {
        auto start = high_resolution_clock::now();
        int result = searchFunc(data, target);
        auto end = high_resolution_clock::now();
        
        auto duration = duration_cast<nanoseconds>(end - start);
        cout << name << ": " << duration.count() << " nanoseconds";
        
        if (result != -1) {
            cout << " (found at index " << result << ")";
        } else {
            cout << " (not found)";
        }
        cout << endl;
    }
};

int main() {
    cout << "=== Sorting and Searching Algorithms Demo ===" << endl;
    
    // Demo array
    vector<int> arr = {64, 34, 25, 12, 22, 11, 90, 88, 76, 50, 42};
    cout << "\nOriginal array: ";
    for (int num : arr) {
        cout << num << " ";
    }
    cout << endl;
    
    // Demonstrate different sorting algorithms
    vector<int> temp;
    
    // Bubble Sort
    temp = arr;
    BubbleSort::bubbleSortOptimized(temp);
    cout << "\nBubble Sort: ";
    for (int num : temp) {
        cout << num << " ";
    }
    cout << endl;
    
    // Quick Sort
    temp = arr;
    QuickSort::quickSort(temp);
    cout << "Quick Sort: ";
    for (int num : temp) {
        cout << num << " ";
    }
    cout << endl;
    
    // Merge Sort
    temp = arr;
    MergeSort::mergeSort(temp);
    cout << "Merge Sort: ";
    for (int num : temp) {
        cout << num << " ";
    }
    cout << endl;
    
    // Heap Sort
    temp = arr;
    HeapSort::heapSort(temp);
    cout << "Heap Sort: ";
    for (int num : temp) {
        cout << num << " ";
    }
    cout << endl;
    
    // Searching demonstrations
    vector<int> sortedArr = {11, 12, 22, 25, 34, 42, 50, 64, 76, 88, 90};
    int target = 42;
    
    cout << "\n=== Searching Demonstrations ===" << endl;
    cout << "Sorted array: ";
    for (int num : sortedArr) {
        cout << num << " ";
    }
    cout << "\nSearching for: " << target << endl;
    
    // Linear Search
    int result = LinearSearch::linearSearch(sortedArr, target);
    cout << "Linear Search: " << (result != -1 ? "Found at index " + to_string(result) : "Not found") << endl;
    
    // Binary Search
    result = BinarySearch::binarySearch(sortedArr, target);
    cout << "Binary Search: " << (result != -1 ? "Found at index " + to_string(result) : "Not found") << endl;
    
    // Jump Search
    result = JumpSearch::jumpSearch(sortedArr, target);
    cout << "Jump Search: " << (result != -1 ? "Found at index " + to_string(result) : "Not found") << endl;
    
    // Advanced applications
    cout << "\n=== Advanced Applications ===" << endl;
    
    // Find kth largest
    vector<int> nums = {3, 2, 1, 5, 6, 4};
    int k = 2;
    int kthLargest = SortingSearchingApplications::findKthLargest(nums, k);
    cout << k << "nd largest element: " << kthLargest << endl;
    
    // Search in rotated array
    vector<int> rotated = {4, 5, 6, 7, 0, 1, 2};
    target = 0;
    result = SortingSearchingApplications::searchRotated(rotated, target);
    cout << "Search " << target << " in rotated array: " << (result != -1 ? "Found at index " + to_string(result) : "Not found") << endl;
    
    // Performance benchmarking
    cout << "\n=== Performance Benchmarking ===" << endl;
    SortingBenchmark::benchmarkSortingAlgorithms();
    SortingBenchmark::benchmarkSearchingAlgorithms();
    
    return 0;
}
```

---

## Summary

### Key Concepts Covered
1. **Simple Sorting**: Bubble, Selection, Insertion Sort
2. **Efficient Sorting**: Merge, Quick, Heap Sort
3. **Specialized Sorting**: Counting, Radix, Bucket Sort
4. **Basic Searching**: Linear and Binary Search
5. **Advanced Searching**: Ternary, Exponential, Interpolation, Jump Search
6. **Applications**: Real-world problems and optimizations

### Time and Space Complexities

#### Sorting Algorithms
| Algorithm | Best Case | Average Case | Worst Case | Space | Stable |
|-----------|-----------|--------------|------------|-------|---------|
| Bubble Sort | O(n) | O(n²) | O(n²) | O(1) | Yes |
| Selection Sort | O(n²) | O(n²) | O(n²) | O(1) | No |
| Insertion Sort | O(n) | O(n²) | O(n²) | O(1) | Yes |
| Merge Sort | O(n log n) | O(n log n) | O(n log n) | O(n) | Yes |
| Quick Sort | O(n log n) | O(n log n) | O(n²) | O(log n) | No |
| Heap Sort | O(n log n) | O(n log n) | O(n log n) | O(1) | No |
| Counting Sort | O(n + k) | O(n + k) | O(n + k) | O(k) | Yes |
| Radix Sort | O(d(n + k)) | O(d(n + k)) | O(d(n + k)) | O(n + k) | Yes |
| Bucket Sort | O(n + k) | O(n + k) | O(n²) | O(n) | Yes |

#### Searching Algorithms
| Algorithm | Best Case | Average Case | Worst Case | Space |
|-----------|-----------|--------------|------------|-------|
| Linear Search | O(1) | O(n) | O(n) | O(1) |
| Binary Search | O(1) | O(log n) | O(log n) | O(1) |
| Jump Search | O(1) | O(√n) | O(√n) | O(1) |
| Interpolation Search | O(1) | O(log log n) | O(n) | O(1) |
| Exponential Search | O(1) | O(log n) | O(log n) | O(1) |

### When to Use Different Algorithms

#### Sorting
- **Small datasets (n < 50)**: Insertion Sort
- **Nearly sorted data**: Insertion Sort, Bubble Sort
- **Memory constrained**: Heap Sort, Selection Sort
- **Stability required**: Merge Sort, Insertion Sort
- **Average case performance**: Quick Sort
- **Guaranteed O(n log n)**: Merge Sort, Heap Sort
- **Integer data with small range**: Counting Sort, Radix Sort
- **Uniformly distributed data**: Bucket Sort

#### Searching
- **Unsorted data**: Linear Search
- **Sorted data**: Binary Search
- **Large sorted arrays**: Binary Search, Jump Search
- **Uniformly distributed sorted data**: Interpolation Search
- **Unknown size arrays**: Exponential Search
- **Unimodal functions**: Ternary Search

### Next Steps
Proceed to **Part 9: Dynamic Programming** to learn about solving complex problems by breaking them into simpler subproblems.

### Additional Resources
- [Sorting Algorithm Animations](https://www.toptal.com/developers/sorting-algorithms)
- [Algorithm Visualizer](https://algorithm-visualizer.org/)
- [Sorting and Searching Practice](https://leetcode.com/tag/sorting/)
- [Big O Cheat Sheet](https://www.bigocheatsheet.com/)