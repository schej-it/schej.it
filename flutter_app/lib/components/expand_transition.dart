import 'package:flutter/material.dart';

class ExpandTransition extends StatelessWidget {
  final bool visible;
  final Widget child;
  final Axis axis;
  final Duration duration;

  const ExpandTransition({
    Key? key,
    required this.child,
    required this.visible,
    this.axis = Axis.vertical,
    this.duration = const Duration(milliseconds: 300),
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AnimatedSwitcher(
      duration: duration,
      switchInCurve: Curves.easeInOut,
      switchOutCurve: Curves.easeInOut,
      transitionBuilder: (child, animation) => SizeTransition(
        axis: axis,
        sizeFactor: animation,
        axisAlignment: 1,
        child: child,
      ),
      child: visible ? child : const SizedBox(width: 0, height: 0),
    );
  }
}
