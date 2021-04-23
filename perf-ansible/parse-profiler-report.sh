awk -F~ '                   ##Setting field separator as tilde here.
{
  val=$2;                   ##Creating a variable named val whose value is 2nd field of current line.
  $2="~@";                   ##Setting value of 2nd field as @ here to keep all lines same(to create index for array a).
  a[$0]+=val                ##Creating an array named a whose index is the current line and its value is the new sum
}
!b[$0]++{                   ##Checking if array b, whose index is the current line, has a value of NULL; if so do following.
  c[++count]=$0}            ##Creating an array named c whose index is variable count increasing value with 1 and value is line.
END{                        ##Starting END block of awk code here.
  for(i=1;i<=count;i++){    ##Starting a for loop whose value starts from 1 to till value of count variable.
     sub("@",a[c[i]],c[i]); ##Substituting @ in value of array c(which is actually lines value) with value of SUMMED $2.
     print c[i]}            ##Printing newly value of array c where $2 is now replaced with its actual value.
}' OFS=\~ <(cat - | sed 's/\(.*\) -\+ \(.*\)s/\1~\2/') | column -s~ -t
