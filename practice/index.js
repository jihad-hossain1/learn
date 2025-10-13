/*
Problem Understanding:
Reverse the order of elements in a string or array.

Key Concepts:
- String/Array manipulation
- Two-pointer technique
- In-place vs new collection
*/


// Method 1: built in method
function reverseStringBuiltIn(str){
    return str.split("").reverse().join('');
}

function reverseArrayBuiltIn(arr){
    return [...arr].reverse();
}

function reverseStringTwoPointers(str){
    const chars = str.split('');
    let left = 0, right = chars.length - 1;

    while(left < right){
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

console.log(result)

console.log(reverseResult)