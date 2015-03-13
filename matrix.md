### 3D Matrix transformation decomposition ###
```
M11	M12	M13	M14
M21	M22	M23	M24
M31	M32	M33	M34
```

### Identity ###
```
1	0	0	0
0	1	0	0
0	0	1	0
```

### Translation Vector ###
```
vt = (M13 M23 M34)
```

### Scaling coefficient ###
```
sx = sqrt(M11² + M12² + M13²);
sy = sqrt(M21² + M22² + M23²);
sz = sqrt(M31² + M32² + M33²);
```

### Rotation Matrix ###
```
M11/sx	M12/sx	M13/sx	0
M21/sy	M22/sy	M23/sy	0
M31/sz	M32/sz	M33/sz	0
```

### 2D Matrix transformation decomposition ###
```
M11	M12	M13
M21	M22	M23
```

### Identity ###
```
1	0	0
0	1	0
```

### Translation Vector ###
```
vt = (M13 M23)
```

### Scaling coefficient ###
```
sx = sqrt(M11² + M12²);
sy = sqrt(M21² + M22²);
sz = sqrt(M31² + M32²);
```

### Rotation Matrix ###
```
M11/sx	M12/sx	0
M21/sy	M22/sy	0
```