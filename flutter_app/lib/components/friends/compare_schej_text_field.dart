import 'package:flutter/material.dart';
import 'package:flutter_app/components/friends/compare_schej_controller.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/api.dart';
import 'package:provider/provider.dart';

// Widget containing chips of all the people added to the current schej comparison
// as well as a text field at the right to search for more people
class CompareSchejTextField extends StatefulWidget {
  final CompareSchejController controller;
  final FocusNode? focusNode;
  final TextEditingController? textEditingController;
  final ScrollController? scrollController;
  final ValueChanged<String>? onSubmitted;

  const CompareSchejTextField({
    Key? key,
    required this.controller,
    this.focusNode,
    this.textEditingController,
    this.scrollController,
    this.onSubmitted,
  }) : super(key: key);

  @override
  State<CompareSchejTextField> createState() => _CompareSchejTextFieldState();
}

class _CompareSchejTextFieldState extends State<CompareSchejTextField> {
  // Variables
  bool _focused = false;
  int _oldLength = 0;

  @override
  void initState() {
    super.initState();

    if (widget.focusNode != null) widget.focusNode!.addListener(_setFocused);
    widget.controller.addListener(
        _clearTextIfUserAdded, [CompareSchejControllerProperties.userIds]);
  }

  @override
  void dispose() {
    if (widget.focusNode != null) widget.focusNode!.removeListener(_setFocused);
    widget.controller.removeListener(
        _clearTextIfUserAdded, [CompareSchejControllerProperties.userIds]);

    super.dispose();
  }

  void _setFocused() {
    if (widget.focusNode != null) {
      setState(() {
        _focused = widget.focusNode!.hasFocus;
      });
    }
  }

  void _clearTextIfUserAdded() {
    if (widget.controller.userIds.length > _oldLength) {
      if (widget.textEditingController != null) {
        widget.textEditingController!.clear();
      }
    }
  }

  // Set the given user as active when tag is tapped
  void _onTagTapped(String? userId) {
    if (userId != null) {
      if (widget.controller.activeUserId == userId) {
        widget.controller.activeUserId = null;
      } else {
        widget.controller.activeUserId = userId;
      }
    } else {
      // The tapped tag is the authUser
      final api = context.read<ApiService>();
      final authUserId = api.authUser!.id;
      if (widget.controller.activeUserId == authUserId) {
        widget.controller.activeUserId = null;
      } else {
        widget.controller.activeUserId = authUserId;
      }
    }
  }

  // Returns whether the given user's tag should display as active
  bool _isUserActive(String? userId) {
    if (widget.controller.activeUserId == null) return true;

    if (userId != null) {
      return widget.controller.activeUserId == userId;
    } else {
      // The tapped tag is the authUser
      final api = context.read<ApiService>();
      final authUserId = api.authUser!.id;
      return widget.controller.activeUserId == authUserId;
    }
  }

  @override
  Widget build(BuildContext context) {
    // Keep track of old length of addedUsers to detect when it increases
    _oldLength = widget.controller.userIds.length;

    return GestureDetector(
      onTap: () {}, // Needed so open container doesn't intercept this event
      child: Container(
        height: 45,
        width: double.infinity,
        decoration: BoxDecoration(
          color: _focused ? SchejColors.white : SchejColors.offWhite,
          borderRadius: SchejConstants.borderRadius,
          border: Border.all(width: 1, color: SchejColors.lightGray),
        ),
        child: Consumer<CompareSchejController>(
          builder: (context, controller, child) => SingleChildScrollView(
            controller: widget.scrollController,
            scrollDirection: Axis.horizontal,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.start,
              mainAxisSize: MainAxisSize.max,
              children: [
                ..._getTags(controller),
                _buildTextField(),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildTextField() {
    return ConstrainedBox(
      constraints: const BoxConstraints(maxWidth: 255),
      child: TextField(
        focusNode: widget.focusNode,
        controller: widget.textEditingController,
        onSubmitted: widget.onSubmitted,
        autocorrect: false,
        decoration: const InputDecoration(
          isDense: true,
          hintText: 'Add a friend...',
          contentPadding: EdgeInsets.symmetric(
            horizontal: 12,
            vertical: 5,
          ),
          enabledBorder: InputBorder.none,
          focusedBorder: InputBorder.none,
        ),
        style: SchejFonts.subtitle,
      ),
    );
  }

  // Returns a list of all the tags to display
  List<Widget> _getTags(CompareSchejController controller) {
    final api = context.read<ApiService>();
    List<Widget> tags = <Widget>[];

    if (controller.includeSelf) {
      tags.add(_buildTag('Me', allowRemove: false));
    }

    for (String userId in controller.userIds) {
      final friend = api.getFriendById(userId);
      tags.add(_buildTag(friend!.fullName, userId: userId));
    }
    return tags;
  }

  // Builds an individual tag
  Widget _buildTag(
    String name, {
    bool allowRemove = true,
    String? userId,
  }) {
    Color textColor;
    Color containerColor;
    if (_isUserActive(userId)) {
      textColor = SchejColors.white;
      containerColor = SchejColors.darkGreen;
    } else {
      textColor = SchejColors.darkGreen;
      containerColor = SchejColors.white;
    }

    return Container(
      decoration: BoxDecoration(
        borderRadius: SchejConstants.borderRadius,
        color: containerColor,
        border: Border.all(color: SchejColors.darkGreen),
      ),
      margin: const EdgeInsets.only(left: 5),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          borderRadius: SchejConstants.borderRadius,
          onTap: () {
            _onTagTapped(userId);
          },
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
            child: Row(
              children: [
                Text(name,
                    style: SchejFonts.subtitle.copyWith(color: textColor)),
                if (allowRemove) ...[
                  const SizedBox(width: 4),
                  InkWell(
                    child: Icon(
                      Icons.cancel,
                      size: 14,
                      color: textColor,
                    ),
                    onTap: () {
                      if (userId != null) {
                        widget.controller.removeUserId(userId);
                      }
                    },
                  ),
                ]
              ],
            ),
          ),
        ),
      ),
    );
  }
}
