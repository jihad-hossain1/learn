/*
Problem Understanding:
Reverse the order of elements in a string or array.

Key Concepts:
- String/Array manipulation
- Two-pointer technique
- In-place vs new collection
*/


// Method 1: built in method
function reverseStringBuiltIn(str) {
    return str.split("").reverse().join('');
}

function reverseArrayBuiltIn(arr) {
    return [...arr].reverse();
}

function reverseStringTwoPointers(str) {
    const chars = str.split('');
    let left = 0, right = chars.length - 1;

    while (left < right) {
        [chars[left], chars[right]] = [chars[right], chars[left]];

        left++;

        right--;
    }

    return chars.join('')
}


// const result = reverseStringBuiltIn("I Love Vegetable")
// const result = reverseArrayBuiltIn([{a: 1},{a:3},{a:9}])
const result = reverseStringTwoPointers("abc xyz abc")
const reverseResult = reverseStringTwoPointers(result)



// console.log(result)

// console.log(reverseResult)

function ownPointerReverse(str) {
    const chars = str.split("")
    let left = 0, right = str.length - 1;

    while (left < right) {
        [chars[left], chars[right]] = [chars[right], chars[left]]

        left++;
        right--;
    }
    return chars.join('')

}

// console.log("own-pointer : ",ownPointerReverse('abc'))



/*
2. Find Max/Min in Array
Problem Understanding:
Find the maximum and minimum values in an array.

Key Concepts:

Array traversal
Comparison operations
Edge cases handling
*/

function findMaxMinValueInList(list) {
     const startTime = performance.now(); // Start timing

    if (list.length == 0) return { max: undefined, min: undefined };

    return {
        max: Math.max(...list),
        min: Math.min(...list)
    }
}

// console.log(findMaxMinValueInList([2,4,52,24,42,44,53,2,5,6,-1,-10,44,3]))

const arr1 = [2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3,2,4,34,5,32,5,77,74,52,24,42,44,53,2,5,6,-1,-10,44,3]
const arr2 = [...arr1];
const arr4 = [...arr1,...arr2];
const arr5 = [...arr1, ...arr4, ...arr2];
const arr6 = [...arr1, ...arr4, ...arr2, ...arr5];
const arr7 = [...arr1, ...arr2, ...arr4, ...arr5, ...arr6];

const arr3 = [...arr1, ...arr2, ...arr4, ...arr5, ...arr6, ...arr7]

console.log("arr3.length", arr3.length);

function findMaxMinValueInList(list) {
    const startTime = performance.now(); // Start timing
    
    if (list.length == 0) {
        const endTime = performance.now();
        console.log(`Execution time: ${(endTime - startTime).toFixed(4)} ms`);
        return { max: undefined, min: undefined };
    }

    const result = {
        max: Math.max(...list),
        min: Math.min(...list)
    };

    const endTime = performance.now();
    console.log(`Built-in Execution time: ${(endTime - startTime).toFixed(4)} ms`);
    
    return result;
}
// console.log(`Built result`, findMaxMinValueInList(arr3));

console.log(` \n`)

function findMaxMinLinear(list) {
    const startTime = performance.now(); // Start timing
    
    if (list.length === 0) {
        const endTime = performance.now();
        console.log(`Execution time: ${(endTime - startTime).toFixed(4)} ms`);
        return { max: undefined, min: undefined };
    }

    let max = list[0]; // Fixed: changed 'arr' to 'list'
    let min = list[0]; // Fixed: changed 'arr' to 'list'

    for (let i = 1; i < list.length; i++) { // Fixed: changed 'arr' to 'list'
        if (list[i] > max) max = list[i];
        if (list[i] < min) min = list[i];
    }

    const endTime = performance.now();
    console.log(`Linear Search Execution time: ${(endTime - startTime).toFixed(4)} ms`);
    
    return { max, min };
}

// Test the function
// console.log("Linear result", findMaxMinLinear(arr3));


function linerSearch(list){
    const startTime = performance.now();

    if(list.length === 0){
        const endTime = performance.now();
        console.log(`execute t : ${(endTime - startTime).toFixed(4)} ms`)
        return {max: undefined, min: undefined}
    }

    let max = list[0];
    let min = list[0];

    for(let i = 1; i < list.length ; i ++ ){
        if ( list[i] > max ) max = list[i];
        if ( list[i] < min ) min = list[i];
    }

    const endTime = performance.now();
    console.log(`ext t : ${(endTime - startTime).toFixed(4)} ms `)

    return {max, min};

}

console.log('ls ', findMaxMinLinear(arr3))





