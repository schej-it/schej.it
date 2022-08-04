import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/components/friends/compare_schej_card.dart';
import 'package:flutter_app/components/friends/compare_schej_text_field.dart';
import 'package:flutter_app/components/friends/compare_schej_text_field_controller.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/api.dart';
import 'package:flutter_app/models/user.dart';
import 'package:provider/provider.dart';

class CompareSchejDialog extends StatefulWidget {
  final CompareSchejTextFieldController controller;
  final VoidCallback onClose;

  // TODO: need to automatically scroll textfield to end when dialog is opened
  const CompareSchejDialog({
    Key? key,
    required this.controller,
    required this.onClose,
  }) : super(key: key);

  @override
  State<CompareSchejDialog> createState() => _CompareSchejDialogState();
}

class _CompareSchejDialogState extends State<CompareSchejDialog> {
  // Controllers
  final _focusNode = FocusNode();
  final _textEditingController = TextEditingController();

  // Variables
  String _query = '';

  @override
  initState() {
    super.initState();

    // When widget is built, automatically focus text field
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _focusNode.requestFocus();
    });

    _textEditingController.addListener(_textChanged);
  }

  @override
  void dispose() {
    _textEditingController.removeListener(_textChanged);
    _textEditingController.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  void _textChanged() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      setState(() => _query = _textEditingController.text);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SchejAppBar(
        isRoot: true,
        underline: false,
        centerTitle: true,
        titleString: 'Add people',
        actions: [
          IconButton(onPressed: widget.onClose, icon: const Icon(Icons.close)),
        ],
      ),
      body: Container(
        color: SchejColors.white,
        padding: SchejConstants.pagePadding,
        child: Column(
          children: [
            CompareSchejTextField(
              focusNode: _focusNode,
              textEditingController: _textEditingController,
              controller: widget.controller,
            ),
            _buildCheckbox(),
            Expanded(child: _buildResults()),
          ],
        ),
      ),
    );
  }

  Widget _buildCheckbox() {
    return Padding(
      padding: const EdgeInsets.only(bottom: 15),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisAlignment: MainAxisAlignment.end,
        children: [
          const Text('Include myself', style: SchejFonts.subtitle),
          Consumer<CompareSchejTextFieldController>(
            builder: (context, controller, child) => Checkbox(
              visualDensity: const VisualDensity(
                vertical: VisualDensity.minimumDensity,
                horizontal: VisualDensity.minimumDensity,
              ),
              value: controller.includeSelf,
              onChanged: (value) => controller.includeSelf = value!,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildResults() {
    return Consumer<ApiService>(builder: (context, api, child) {
      List<User> results = api.getFriendsByQuery(_query);
      return ListView.builder(
        itemCount: results.length,
        itemBuilder: (context, index) {
          final friend = results[index];
          return Padding(
            key: ValueKey(friend.id),
            padding: const EdgeInsets.only(bottom: 10),
            child: Consumer<CompareSchejTextFieldController>(
              builder: (context, controller, child) => CompareSchejCard(
                name: friend.fullName,
                picture: friend.picture,
                added: widget.controller.userIds.contains(friend.id),
                onToggle: (bool value) {
                  if (value) {
                    widget.controller.addUserId(friend.id);
                  } else {
                    widget.controller.removeUserId(friend.id);
                  }
                },
              ),
            ),
          );
        },
      );
    });
  }
}
