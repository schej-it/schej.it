import 'package:animations/animations.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/components/calendar/calendar.dart';
import 'package:flutter_app/components/calendar/calendar_view_selector.dart';
import 'package:flutter_app/components/friends/compare_schej_text_field.dart';
import 'package:flutter_app/components/friends/compare_schej_controller.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/models/api.dart';
import 'package:flutter_app/pages/friends/compare_schej_dialog.dart';
import 'package:flutter_app/utils.dart';
import 'package:linked_scroll_controller/linked_scroll_controller.dart';
import 'package:provider/provider.dart';

class CompareSchejPage extends StatefulWidget {
  final String friendId;
  final bool initialIncludeSelf;

  const CompareSchejPage({
    Key? key,
    required this.friendId,
    this.initialIncludeSelf = false,
  }) : super(key: key);

  @override
  State<CompareSchejPage> createState() => _CompareSchejPageState();
}

class _CompareSchejPageState extends State<CompareSchejPage> {
  // Controllers
  late final CompareSchejController _compareSchejController;
  late final LinkedScrollControllerGroup _controllers;
  late final ScrollController _textFieldScrollController;
  late final ScrollController _dialogScrollController;

  // Calendar variables
  int _daysVisible = 3;
  DateTime _selectedDay = getDateWithTime(DateTime.now(), 0);
  bool _showAvailability = false;

  @override
  void initState() {
    super.initState();

    _compareSchejController = CompareSchejController(
      initialUserIds: <String>{widget.friendId},
      initialActiveUserId: widget.friendId,
      initialIncludeSelf: widget.initialIncludeSelf,
    );
    _compareSchejController.addListener(_setActiveUserId, [
      CompareSchejControllerProperties.userIds,
      CompareSchejControllerProperties.includeSelf,
    ]);

    _controllers = LinkedScrollControllerGroup();
    _textFieldScrollController = _controllers.addAndGet();
    _dialogScrollController = _controllers.addAndGet();
  }

  @override
  void dispose() {
    _compareSchejController.removeListener(_setActiveUserId);
    _compareSchejController.dispose();
    _textFieldScrollController.dispose();
    _dialogScrollController.dispose();

    super.dispose();
  }

  // Sets the active user id when the compareSchejController changes
  void _setActiveUserId() {
    int numUsers = _compareSchejController.userIds.length;
    if (_compareSchejController.includeSelf) numUsers++;

    if (numUsers == 2) {
      _compareSchejController.activeUserId =
          _compareSchejController.userIds.first;
    } else {
      _compareSchejController.activeUserId = null;
    }

    if (numUsers > 2 && !_showAvailability) {
      setState(() {
        _showAvailability = true;
      });
    } else if (numUsers <= 2 && _showAvailability) {
      setState(() {
        _showAvailability = false;
      });
    }
  }

  Set<String> _getUserIds(ApiService api, CompareSchejController controller) {
    if (!controller.includeSelf) return controller.userIds;

    return Set.from(controller.userIds)..add(api.authUser!.id);
  }

  @override
  Widget build(BuildContext context) {
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
      body: ChangeNotifierProvider.value(
        value: _compareSchejController,
        child: Container(
          color: SchejColors.white,
          child: Stack(
            children: [
              Positioned.fill(
                child: Consumer2<ApiService, CompareSchejController>(
                  builder: (context, api, controller, child) => Calendar(
                    mode: CalendarMode.compare,
                    userIds: _getUserIds(api, controller),
                    daysVisible: _daysVisible,
                    selectedDay: _selectedDay,
                    onDaySelected: (selectedDay) => setState(() {
                      _selectedDay = selectedDay;
                    }),
                    showAvatars: true,
                    showAvailability: _showAvailability,
                    activeUserId: controller.activeUserId,
                  ),
                ),
              ),
              _buildTextField(),
            ],
          ),
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
            return ChangeNotifierProvider.value(
              value: _compareSchejController,
              child: FocusScope(
                child: Focus(
                  onFocusChange: (focus) {
                    if (focus) {
                      openContainer();
                    }
                  },
                  child: CompareSchejTextField(
                    controller: _compareSchejController,
                    scrollController: _textFieldScrollController,
                  ),
                ),
              ),
            );
          },
          openBuilder: (context, closeContainer) {
            return ChangeNotifierProvider.value(
              value: _compareSchejController,
              child: CompareSchejDialog(
                controller: _compareSchejController,
                scrollController: _dialogScrollController,
                onClose: closeContainer,
              ),
            );
          },
        ),
      ),
    );
  }
}
