import 'package:flutter/material.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

class CalendarViewSelector extends StatelessWidget {
  final Function(int) onSelected;

  const CalendarViewSelector({
    Key? key,
    required this.onSelected,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return PopupMenuButton(
      icon: const Icon(MdiIcons.calendarBlankOutline),
      splashRadius: 15,
      onSelected: onSelected,
      itemBuilder: (context) => [
        const PopupMenuItem(
          value: 1,
          child: Text('Day'),
        ),
        const PopupMenuItem(
          value: 3,
          child: Text('3 day'),
        ),
        const PopupMenuItem(
          value: 7,
          child: Text('Week'),
        ),
      ],
    );
  }
}
