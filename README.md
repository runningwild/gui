A widgets package for go
========================

This is an experiment in creating a "widgets" package for go.  The
idea is to create a widget library utilizing functions that combine
together multiple widgets to form larger widgets.  This is sort of a
more "functional programming" (and hopefully more go-like) style than
conventional GUI libraries, which tend to be strongly oriented toward
inheritence.

I would like to make the widget library at least somewhat independent
of the rendering back end.  The current back end is based on
websockets and html 5 with a bit of javascript in between.  I don't
think it's a very nice design, but the point is to create a nice API
that isn't dependent on the specific back end.  Not that I've created
this yet, but it's my goal.  I'd like to eventually be able to support
a GTK version (or native windows) of the same application with (close
to?) no change in the source code of the application.

So far, this framework is highly incomplete.  All I have are buttons
and non-user-editable text.  If this README gets out of date, I may
have more widget types implemented.
