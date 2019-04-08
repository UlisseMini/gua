# lua + a bunch of exported go functions

## Example
```lua
local addr = 'localhost:1337'
local conn, err = gua.dial(addr)
if err ~= nil then
  print('error: '..err)
  return
end

print('connected to '..addr)

print 'sending...'
conn.write('foo')
conn.write('bar from lua\n')

print 'reading...'
data, err = conn.read(1024)
if err ~= nil then
  print('error: '..err)
else
  print('got: '..data)
end

print 'closing...'
conn.close()

print 'all done!'
```
