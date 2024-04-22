# MDCalc
A language made to have mathmatical expressions automatically be converted to a list of calculations. (for the exam where the calculations must be explained)  
# Docs
## Syntax
MDCalc renders instructions line by line, starting in 1.mdc, then 2.mdc, 3.mdc and so on.  
Every n.mdc defines a solution for problem n (so 1.mdc for problem 1).  
Every line must start with an instruction character.  
### Instruction Syntax
| Instruction | function |
| - | - |
| \| | Starts a new subproblem, anything before the first of these will not be rendered |
| T [text] | Renders text |
| C [expr] | Evaluates and renders expression, can be used before \| to init variables
| I [image name] | Renders an image |
### Expression Syntax
| Name | Syntax | Function |
| - | - | - |
| Literal | *{number}* or *{varname}* | Represents a literal (already known) number.
| Variable Setter | *({varname} = {expr})* | Sets a variable for use in later expressions and returns the result of the expression, variables also store their unit. |
| Comment | *({expr}:{text})* or *({expr}:{int}:{text})* | Renders comment after expression, and optionally sets how many decimals should be rendered. Comments will split a calculation into multiple lines, unless the resulting line will be "literal = literal" (this is so precision can be set without an extra line). |
| Unit override | *{number}{unit}* or *({expr}){unit}* (can be used after function) | Overrides unit of expression result or literal, when formatted literals will display their units, and the result will display its unit. When compiling MDCalc will ask for unit display names. |
| Operator | *{expr}{op}{expr}* | Applies operator to expressions. When compiling MDCalc will ask for the resulting unit of applying operators to units.  |
| Function | *{function}({expr}, {expr}, ...)* | Applies function to expression(s), resulting unit will always be None. |

*note: it is not necessary to add parathesis around an entire calculation or function parameter for comments and var setters*

### Operators
| Operator | Function |
| - | - |
| * | Multiply
| / | Divide with fraction line
| % | Also divides, just formats without fraction line, for mod use the function
| ^ | Power
| + | Addition
| - | Subtraction

### Functions

*all trigonometry uses degrees*
| Name | Comment |
| - | - |
| pi() | Always returns pi
| e() | Always returns e
| par(x) | Adds parenthesis around formatted output, use in case MDCalc fails to add them
| neg(x) | Negates the input and formats using a '-', this is because I forgot to add negating to the MDCalc syntax, remember to use par() on input if parenthesis are needed 
| floor(x)
| ceil(x)
| abs(x)
| sqrt(x)
| root(r, x)
| log10(x)
| log(b, x)
| sin(x)
| cos(x)
| tan(x)
| atan(x)
| asin(x)
| acos(x)
| mod(a, b)