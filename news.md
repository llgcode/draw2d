# How to use draw2d with Gotour #
[Draw2d with Gotour](Gotour.md)

# draw2d and opengl #
date: 5/15/2011 author: Laurent Le Goff

I'ts now possible to use the common GraphicContext to display 2d figure on opengl. The particularity of this implementationis that it is compatible with old opengl implementation. And either the graphic card the rendering is the same. I use Vextex buffer to draw hspan with the freetype Painter.
You can test the opengl implementation (tested on windows and macosx) http://code.google.com/p/draw2d/source/browse/cmd/draw2dgl.go

If you test this code please tell me about your metrics displayed in console.


# Test draw2d on windows 7 #
date: 2/15/2011 author: Laurent Le Goff

I've recently installed go on a windows 7 machine.
I've tested the possibility of using the wingui and added some support for bitmap.
My first simple test works nice the Go program parse a postscript file  and display it in a window on win7.
![http://draw2d.googlecode.com/svn/wiki/news/tiger_viewer.png](http://draw2d.googlecode.com/svn/wiki/news/tiger_viewer.png)