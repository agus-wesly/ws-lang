// // Block test 
// let a = "global a";
// let b = "global b";
// let c = "global c";
// {
//   let a = "outer a";
//   let b = "outer b";
//   {
//     let a = "inner a";
//     print a;
//     print b;
//     print c;
//   }
//   print a;
//   print b;
//   print c;
// }
// print a;
// print b;
// print c;
// 
// // While test 
// let _a = 0;
// let temp;
// 
// for (let b = 1; _a < 10000; b = temp + b) {
//   print _a;
//   temp = _a;
//   _a = b;
// }


fun count(n) {
  if (n > 1) count(n - 1);
  print n;
}

count(10);

// foo();
// 
// foo(2);
// 
// fun foo(x,y,z)  {
//     return x +y;
// }



