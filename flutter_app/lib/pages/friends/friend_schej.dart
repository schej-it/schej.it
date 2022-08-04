import 'package:animations/animations.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/components/calendar.dart';
import 'package:flutter_app/components/calendar_view_selector.dart';
import 'package:flutter_app/components/friends/compare_schej_text_field.dart';
import 'package:flutter_app/components/friends/compare_schej_text_field_controller.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/models/calendar_event.dart';
import 'package:flutter_app/pages/friends/compare_schej_dialog.dart';
import 'package:flutter_app/utils.dart';
import 'package:provider/provider.dart';

class FriendSchejPage extends StatefulWidget {
  final String friendId;

  const FriendSchejPage({Key? key, required this.friendId}) : super(key: key);

  @override
  State<FriendSchejPage> createState() => _FriendSchejPageState();
}

class _FriendSchejPageState extends State<FriendSchejPage> {
  // Controllers
  late final CompareSchejTextFieldController _compareSchejTextFieldController;

  // Calendar variables
  int _daysVisible = 3;
  DateTime _selectedDay = getDateWithTime(DateTime.now(), 0);

  @override
  void initState() {
    super.initState();

    _compareSchejTextFieldController = CompareSchejTextFieldController(
      initialUserIds: <String>{widget.friendId},
      initialIncludeSelf: true,
    );
  }

  @override
  void dispose() {
    _compareSchejTextFieldController.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final testCalendarEvents = CalendarEvents(
      events: [
        CalendarEvent(
          title: 'Event',
          startDate: getDateWithTime(DateTime.now(), 9.5),
          endDate: getDateWithTime(DateTime.now(), 12),
        ),
        CalendarEvent(
          title: 'Introduction to Failure Analysis',
          startDate: getDateWithTime(DateTime.now(), 13),
          endDate: getDateWithTime(DateTime.now(), 14.5),
        ),
        CalendarEvent(
          title: 'War',
          startDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 1)), 15),
          endDate:
              getDateWithTime(DateTime.now().add(const Duration(days: 1)), 20),
        ),
      ],
    );

    return Scaffold(
      appBar: SchejAppBar(
        centerTitle: true,
        underline: false,
        titleString: 'Compare schej',
        actions: [
          CalendarViewSelector(
            onSelected: (int value) => setState(() {
              _daysVisible = value;
            }),
          ),
        ],
      ),
      body: Container(
        color: SchejColors.white,
        child: Stack(
          children: [
            Positioned.fill(
              child: Calendar(
                calendarEvents: testCalendarEvents,
                daysVisible: _daysVisible,
                selectedDay: _selectedDay,
                onDaySelected: (selectedDay) => setState(() {
                  _selectedDay = selectedDay;
                }),
              ),
            ),
            _buildTextField(),
          ],
        ),
      ),
    );
  }

  Widget _buildTextField() {
    return Positioned(
      bottom: 0,
      width: MediaQuery.of(context).size.width,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 25, horizontal: 18),
        child: OpenContainer(
          closedShape: RoundedRectangleBorder(
            borderRadius: SchejConstants.borderRadius,
          ),
          closedBuilder: (context, openContainer) {
            // TODO: don't wrap this with focus widget because it should be
            // the raw text field that is focused, right now if I tap ANYWHERE
            // in the text field it focuses, even if I tapped a tag
            // Maybe have an onFocus callback??
            return ChangeNotifierProvider.value(
              value: _compareSchejTextFieldController,
              child: FocusScope(
                child: Focus(
                  onFocusChange: (focus) {
                    if (focus) {
                      openContainer();
                    }
                  },
                  child: CompareSchejTextField(
                    controller: _compareSchejTextFieldController,
                  ),
                ),
              ),
            );
          },
          openBuilder: (context, closeContainer) {
            return ChangeNotifierProvider.value(
              value: _compareSchejTextFieldController,
              child: CompareSchejDialog(
                controller: _compareSchejTextFieldController,
                onClose: closeContainer,
              ),
            );
          },
        ),
      ),
    );
  }
}
