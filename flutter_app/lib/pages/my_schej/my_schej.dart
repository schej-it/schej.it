import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/components/calendar.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

class MySchejPage extends StatefulWidget {
  const MySchejPage({Key? key}) : super(key: key);

  @override
  State<MySchejPage> createState() => _MySchejPageState();
}

class _MySchejPageState extends State<MySchejPage> {
  bool _monthSelector = false;
  bool _viewMenu = false;
  bool _eventNamesVisible = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildAppBar(),
      body: Container(
        color: SchejColors.white,
        padding: SchejConstants.pagePadding,
        child: const Calendar(),
      ),
    );
  }

  PreferredSizeWidget _buildAppBar() {
    return SchejAppBar(
      title: Row(
        children: [
          const Text('July'),
          IconButton(
            icon: _monthSelector
                ? const Icon(MdiIcons.chevronUp)
                : const Icon(MdiIcons.chevronDown),
            splashRadius: 15,
            onPressed: () => setState(() {
              _monthSelector = !_monthSelector;
            }),
          ),
        ],
      ),
      actions: [
        IconButton(
          icon: const Icon(MdiIcons.calendarBlankOutline),
          splashRadius: 15,
          onPressed: () => setState(() {
            _viewMenu = !_viewMenu;
          }),
        ),
        IconButton(
          icon: const Icon(MdiIcons.eye),
          splashRadius: 15,
          onPressed: () => setState(() {
            _eventNamesVisible = !_eventNamesVisible;
          }),
        ),
      ],
      underline: false,
      isRoot: true,
    );
  }
}
