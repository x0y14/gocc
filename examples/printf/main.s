.text
.align 2
.global _main

_main:
  stp x29, x30, [sp, #-32]! ; @A sp-=32, push x29, push 30, maybe

  bl _foo                   ; return = 2
  str x0, [sp, #-16]!       ; prepare x8 = 2

  bl _foo                   ; return = 2
  add x0, x0, #10           ; return += 10
  str x0, [sp, #-16]!       ; prepare x9 = 12

  ldr x9, [sp], #16         ; x9 = 12
  ldr x8, [sp], #16         ; x8 = 2
  add x8, x8, x9            ; x8 = x8 + x9
  mov x0, x8                ; x0 = x8 = 14

  ldp x29, x30, [sp], #32   ; @A
  ret                       ; main's return

