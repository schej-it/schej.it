import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/api.dart';
import 'package:provider/provider.dart';

// Widget containing chips of all the people added to the current schej comparison
// as well as a text field at the right to search for more people
class CompareSchejTextField extends StatefulWidget {
  final FocusNode? focusNode;
  final TextEditingController? controller;
  final Set<String> addedUsers;
  final bool includeSelf;
  // final ValueChanged<String>? onChanged;
  final void Function(String userId) onRemoveUser;

  const CompareSchejTextField({
    Key? key,
    this.focusNode,
    this.controller,
    required this.addedUsers,
    this.includeSelf = true,
    // this.onChanged,
    required this.onRemoveUser,
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
  }

  @override
  void dispose() {
    if (widget.focusNode != null) widget.focusNode!.removeListener(setFocused);
    super.dispose();
  }

  @override
  void didUpdateWidget(covariant CompareSchejTextField oldWidget) {
    super.didUpdateWidget(oldWidget);

    // If user was added, clear text field
    if (widget.addedUsers.length > _oldLength) {
      if (widget.controller != null) {
        widget.controller!.clear();
      }
      // if (widget.onChanged != null) {
      //   widget.onChanged!('');
      // }
    }
  }

  void setFocused() {
    if (widget.focusNode != null) {
      setState(() {
        _focused = widget.focusNode!.hasFocus;
      });
    }
  }

  bool areTagsEmpty() {
    return widget.addedUsers.isEmpty && !widget.includeSelf;
  }

  @override
  Widget build(BuildContext context) {
    // Keep track of old length of addedUsers to detect when it increases
    _oldLength = widget.addedUsers.length;

    final textField = ConstrainedBox(
      constraints: const BoxConstraints(maxWidth: 255),
      child: TextField(
        focusNode: widget.focusNode,
        controller: widget.controller,
        // onChanged: widget.onChanged,
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
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          mainAxisSize: MainAxisSize.max,
          children: [
            ..._getTags(),
            textField,
          ],
        ),
      ),
    );
  }

  // Returns a list of all the tags to display
  List<Widget> _getTags() {
    final api = context.read<ApiService>();
    List<Widget> tags = <Widget>[];

    if (widget.includeSelf) {
      tags.add(_buildTag('Me', allowRemove: false));
    }

    for (String userId in widget.addedUsers) {
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
                  widget.onRemoveUser(userId);
                }
              },
            ),
          ]
        ],
      ),
    );
  }
}
