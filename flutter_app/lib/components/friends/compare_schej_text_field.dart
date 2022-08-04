import 'package:flutter/material.dart';
import 'package:flutter_app/components/friends/compare_schej_text_field_controller.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/api.dart';
import 'package:provider/provider.dart';

// Widget containing chips of all the people added to the current schej comparison
// as well as a text field at the right to search for more people
class CompareSchejTextField extends StatefulWidget {
  final FocusNode? focusNode;
  final TextEditingController? textEditingController;
  final CompareSchejTextFieldController controller;

  const CompareSchejTextField({
    Key? key,
    this.focusNode,
    this.textEditingController,
    required this.controller,
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

    if (widget.focusNode != null) widget.focusNode!.addListener(setFocused);
    widget.controller.addListener(_clearTextIfUserAdded);
  }

  @override
  void dispose() {
    if (widget.focusNode != null) widget.focusNode!.removeListener(setFocused);
    widget.controller.removeListener(_clearTextIfUserAdded);

    super.dispose();
  }

  void setFocused() {
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

  @override
  Widget build(BuildContext context) {
    // Keep track of old length of addedUsers to detect when it increases
    _oldLength = widget.controller.userIds.length;

    final textField = ConstrainedBox(
      constraints: const BoxConstraints(maxWidth: 255),
      child: TextField(
        focusNode: widget.focusNode,
        controller: widget.textEditingController,
        autocorrect: false,
        decoration: const InputDecoration(
          isDense: true,
          hintText: 'Type a name...',
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

    return Container(
      height: 45,
      width: double.infinity,
      decoration: BoxDecoration(
        color: _focused ? SchejColors.white : SchejColors.offWhite,
        borderRadius: SchejConstants.borderRadius,
        border: Border.all(width: 1, color: SchejColors.lightGray),
      ),
      child: Consumer<CompareSchejTextFieldController>(
        builder: (context, controller, child) => SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          child: Row(
            mainAxisAlignment: MainAxisAlignment.start,
            mainAxisSize: MainAxisSize.max,
            children: [
              ..._getTags(controller),
              textField,
            ],
          ),
        ),
      ),
    );
  }

  // Returns a list of all the tags to display
  List<Widget> _getTags(CompareSchejTextFieldController controller) {
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
  Widget _buildTag(String name, {bool allowRemove = true, String? userId}) {
    return Container(
      decoration: BoxDecoration(
        borderRadius: SchejConstants.borderRadius,
        color: SchejColors.green,
      ),
      margin: const EdgeInsets.only(left: 5),
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
      child: Row(
        children: [
          InkWell(
            child: Text(name,
                style: SchejFonts.subtitle.copyWith(color: SchejColors.white)),
            onTap: () {
              // print('$tag selected!');
            },
          ),
          if (allowRemove) ...[
            const SizedBox(width: 4),
            InkWell(
              child: const Icon(
                Icons.cancel,
                size: 14,
                color: SchejColors.white,
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
    );
  }
}
