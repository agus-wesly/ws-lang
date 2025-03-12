let a = "global a";
let b = "global b";
{
  let a = "outer a";
  let b = "outer b";
  {
    let a = "inner a";
    print a;
    print b;
    print c;
    a = "inner a new";
  }
  print a;
  print b;
  print c;
}
print a;
print b;
print c;
